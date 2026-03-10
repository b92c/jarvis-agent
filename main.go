package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "time/tzdata"

	"jarvis-agent/agent"
	"jarvis-agent/config"
	"jarvis-agent/services"
	"jarvis-agent/tools"
)

var saoPauloLocation *time.Location

func init() {
	saoPauloLocation = loadSaoPauloLocation()
}

func loadSaoPauloLocation() *time.Location {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err == nil {
		return loc
	}

	log.Printf("⚠️  Erro ao carregar timezone America/Sao_Paulo: %v. Usando offset fixo UTC-3.", err)
	return time.FixedZone("America/Sao_Paulo", -3*60*60)
}

func main() {
	runOnce := flag.Bool("once", false, "Executa o relatório uma vez e encerra")
	flag.Parse()

	log.Println("🤖 Jarvis Agent Starting...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	log.Printf("📂 Notas: %s", cfg.TodayNotesPath)
	log.Printf("📄 Report: %s", cfg.ReportOutputPath)
	log.Printf("🧠 Modelo: %s", cfg.GeminiModel)

	reader := tools.NewFileReader(cfg)
	writer := tools.NewFileWriter(cfg)
	emailSvc := services.NewEmailService(cfg)

	jarvisAgent, err := agent.NewJarvisAgent(cfg, reader, writer)
	if err != nil {
		log.Fatalf("Erro ao criar agente: %v", err)
	}

	if *runOnce {
		executeReport(jarvisAgent, emailSvc)
		return
	}

	log.Println("✅ Jarvis Agent inicializado com sucesso!")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\n🛑 Encerrando Jarvis...")
		cancel()
	}()

	runSchedulerWithTicker(ctx, jarvisAgent, emailSvc, cfg.StateFilePath)
}

type SchedulerState struct {
	LastExecutionDate string `json:"last_execution_date"`
	Executed12h       bool   `json:"executed_12h"`
	Executed18h       bool   `json:"executed_18h"`
	LastMonthExecuted string `json:"last_month_executed"`
}

func loadState(stateFilePath string) (*SchedulerState, error) {
	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &SchedulerState{}, nil
		}
		return nil, fmt.Errorf("erro ao ler arquivo de estado: %w", err)
	}

	var state SchedulerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo de estado: %w", err)
	}

	return &state, nil
}

func saveState(stateFilePath string, state *SchedulerState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar estado: %w", err)
	}

	dir := stateFilePath[:len(stateFilePath)-len(filepath.Base(stateFilePath))]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de estado: %w", err)
	}

	if err := os.WriteFile(stateFilePath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar arquivo de estado: %w", err)
	}

	return nil
}

func runSchedulerWithTicker(ctx context.Context, jarvisAgent *agent.JarvisAgent, emailSvc *services.EmailService, stateFilePath string) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	log.Println("⏰ Scheduler ativo: executará às 12:00 e 18:00 (diário) e dia 28 às 18:00 (mensal)")
	log.Println("⏰ Verificação a cada 10 minutos")

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler encerrado")
			return
		case <-ticker.C:
			now := time.Now()
			inSaoPaulo := now.In(saoPauloLocation)
			today := inSaoPaulo.Format("2006-01-02")
			currentMonth := inSaoPaulo.Format("2006-01")

			state, err := loadState(stateFilePath)
			if err != nil {
				log.Printf("⚠️ Erro ao carregar estado: %v", err)
				continue
			}

			if state.LastExecutionDate != today {
				state.Executed12h = false
				state.Executed18h = false
			}

			if state.LastMonthExecuted != currentMonth {
				state.LastMonthExecuted = ""
			}

			hour := inSaoPaulo.Hour()
			minute := inSaoPaulo.Minute()
			day := inSaoPaulo.Day()

			if hour == 12 && minute < 10 && !state.Executed12h {
				log.Println("🕐 Executando relatório das 12:00...")
				executeReport(jarvisAgent, emailSvc)
				state.Executed12h = true
				state.LastExecutionDate = today
				if err := saveState(stateFilePath, state); err != nil {
					log.Printf("⚠️ Erro ao salvar estado: %v", err)
				}
				continue
			}

			if hour == 18 && minute < 10 && !state.Executed18h {
				log.Println("🕕 Executando relatório das 18:00...")
				executeReport(jarvisAgent, emailSvc)
				state.Executed18h = true
				state.LastExecutionDate = today
				if err := saveState(stateFilePath, state); err != nil {
					log.Printf("⚠️ Erro ao salvar estado: %v", err)
				}
				continue
			}

			if day == 28 && hour == 18 && minute < 10 && state.LastMonthExecuted != currentMonth {
				log.Println("🗓️ Executando relatório mensal (dia 28)...")
				executeMonthlyReport(jarvisAgent, emailSvc)
				state.LastMonthExecuted = currentMonth
				if err := saveState(stateFilePath, state); err != nil {
					log.Printf("⚠️ Erro ao salvar estado: %v", err)
				}
				continue
			}

			if err := saveState(stateFilePath, state); err != nil {
				log.Printf("⚠️ Erro ao salvar estado: %v", err)
			}
		}
	}
}

func executeReport(jarvisAgent *agent.JarvisAgent, emailSvc *services.EmailService) {
	log.Println("🚀 Executando relatório diário...")

	report, err := jarvisAgent.RunDaily()
	if err != nil {
		log.Printf("❌ Erro ao gerar relatório: %v", err)
		return
	}

	log.Println("📧 Enviando email...")

	if err := emailSvc.SendReport(report); err != nil {
		log.Printf("❌ Erro ao enviar email: %v", err)
		return
	}

	log.Println("✅ Relatório executado e enviado com sucesso!")
}

func executeMonthlyReport(jarvisAgent *agent.JarvisAgent, emailSvc *services.EmailService) {
	log.Println("🚀 Executando relatório mensal...")

	report, err := jarvisAgent.RunMonthly()
	if err != nil {
		log.Printf("❌ Erro ao gerar relatório mensal: %v", err)
		return
	}

	log.Println("📧 Enviando email...")

	if err := emailSvc.SendReport(report); err != nil {
		log.Printf("❌ Erro ao enviar email: %v", err)
		return
	}

	log.Println("✅ Relatório mensal executado e enviado com sucesso!")
}
