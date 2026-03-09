package main

import (
	"context"
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
	log.Println("🤖 Jarvis Agent Starting...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	log.Printf("📂 Notas: %s", cfg.TodayNotesPath)
	log.Printf("📄 Report: %s", cfg.ReportOutputPath)

	reader := tools.NewFileReader(cfg)
	writer := tools.NewFileWriter(cfg)
	emailSvc := services.NewEmailService(cfg)

	jarvisAgent, err := agent.NewJarvisAgent(cfg, reader, writer)
	if err != nil {
		log.Fatalf("Erro ao criar agente: %v", err)
	}

	log.Println("✅ Jarvis Agent inicializado com sucesso!")
	log.Println("⏰ Scheduler ativo: executará às 12:00 e 18:00")

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

func runScheduler(ctx context.Context, jarvisAgent *agent.JarvisAgent, emailSvc *services.EmailService) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler encerrado")
			return
		default:
			now := time.Now()

			nextRun := getNextRunTime(now)
			waitDuration := nextRun.Sub(now)

			log.Printf("⏳ Próxima execução: %s", nextRun.Format("02/01/2006 15:04:05"))
			log.Printf("⏳ Aguardando %v...", waitDuration)

			select {
			case <-ctx.Done():
				return
			case <-time.After(waitDuration):
				executeReport(jarvisAgent, emailSvc)
			}
		}
	}
}

func getNextRunTime(now time.Time) time.Time {
	hour12 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, saoPauloLocation)
	hour18 := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, saoPauloLocation)

	if now.Before(hour12) {
		return hour12
	}
	if now.Before(hour18) {
		return hour18
	}

	nextDay12 := hour12.AddDate(0, 0, 1)
	return nextDay12
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
