package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
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
	log.Println("⏰ Scheduler ativo: executará às 12:00 e 18:00 (diário) e às 18:00 (mensal - último dia do mês)")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\n🛑 Encerrando Jarvis...")
		cancel()
	}()

	runScheduler(ctx, jarvisAgent, emailSvc)
}

type runType int

const (
	runDaily runType = iota
	runMonthly
)

type scheduledRun struct {
	time    time.Time
	runType runType
}

func runScheduler(ctx context.Context, jarvisAgent *agent.JarvisAgent, emailSvc *services.EmailService) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler encerrado")
			return
		default:
			now := time.Now()

			nextRun := getNextRunTime(now)
			waitDuration := nextRun.time.Sub(now)

			runTypeStr := "diário"
			if nextRun.runType == runMonthly {
				runTypeStr = "mensal"
			}
			log.Printf("⏳ Próxima execução: %s (%s)", nextRun.time.Format("02/01/2006 15:04:05"), runTypeStr)
			log.Printf("⏳ Aguardando %v...", waitDuration)

			select {
			case <-ctx.Done():
				return
			case <-time.After(waitDuration):
				if nextRun.runType == runMonthly {
					executeMonthlyReport(jarvisAgent, emailSvc)
				} else {
					executeReport(jarvisAgent, emailSvc)
				}
			}
		}
	}
}

func getNextRunTime(now time.Time) scheduledRun {
	hour12 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, saoPauloLocation)
	hour18 := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, saoPauloLocation)

	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, 18, 0, 0, 0, saoPauloLocation)

	if now.Day() == lastDayOfMonth.Day() && now.After(hour12) {
		return scheduledRun{time: lastDayOfMonth, runType: runMonthly}
	}

	if now.Before(hour12) {
		return scheduledRun{time: hour12, runType: runDaily}
	}
	if now.Before(hour18) {
		return scheduledRun{time: hour18, runType: runDaily}
	}

	nextDay12 := hour12.AddDate(0, 0, 1)
	return scheduledRun{time: nextDay12, runType: runDaily}
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
