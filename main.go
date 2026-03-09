package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jarvis-agent/agent"
	"jarvis-agent/config"
	"jarvis-agent/services"
	"jarvis-agent/tools"
)

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
	log.Println("⏰ Scheduler ativo: executará às 11:50 e 18:00")

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
	hour11 := time.Date(now.Year(), now.Month(), now.Day(), 11, 50, 0, 0, now.Location())
	hour18 := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

	if now.Before(hour11) {
		return hour11
	}
	if now.Before(hour18) {
		return hour18
	}

	nextDay11 := hour11.AddDate(0, 0, 1)
	return nextDay11
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
