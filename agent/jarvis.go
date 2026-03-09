package agent

import (
	"context"
	"fmt"
	"log"

	"jarvis-agent/config"
	"jarvis-agent/tools"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type JarvisAgent struct {
	agentx  agent.Agent
	cfg     *config.Config
	runnerx *runner.Runner
	reader  *tools.FileReader
	writer  *tools.FileWriter
}

func NewJarvisAgent(cfg *config.Config, reader *tools.FileReader, writer *tools.FileWriter) (*JarvisAgent, error) {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-3-pro-preview", &genai.ClientConfig{
		APIKey: cfg.GoogleAPIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar modelo: %w", err)
	}

	jarvisAgent, err := llmagent.New(llmagent.Config{
		Name:        "jarvis_agent",
		Model:       model,
		Description: "Assistente pessoal de Team/Tech Lead que analisa anotações diárias",
		Instruction: getJarvisInstructions(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agente: %w", err)
	}

	sessionService := session.InMemoryService()

	runnerx, err := runner.New(runner.Config{
		Agent:          jarvisAgent,
		SessionService: sessionService,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar runner: %w", err)
	}

	return &JarvisAgent{
		agentx:  jarvisAgent,
		cfg:     cfg,
		runnerx: runnerx,
		reader:  reader,
		writer:  writer,
	}, nil
}

func getJarvisInstructions() string {
	return `
Você é Jarvis, um assistente pessoal de Team/Tech Lead altamente competente.

SUAS RESPONSABILIDADES:
1. Analisar as anotações diárias do usuário em formato Markdown
2. Identificar os PRINCIPAIS PONTOS do dia (reuniões, dailys, rituais ágeis, decisões, etc)
3. Destacar PONTOS DE ATENÇÃO e riscos potenciais
4. Propor AÇÕES PRÁTICAS para problemas mencionados
5. Oferecer SOLUÇÕES para conflitos ou obstáculos descritos

FORMATO DO RELATÓRIO:
Quando gerar o relatório, use o seguinte formato em Markdown:

# 📋 Jarvis Report - [DATA]

## 🎯 Principais Pontos do Dia
- [Ponto 1]
- [Ponto 2]
- [...]

## ⚠️ Pontos de Atenção
- [Atenção 1]
- [Atenção 2]
- [...]

## 💡 Ações Práticas Sugeridas
- [Ação 1]
- [Ação 2]
- [...]

## 🔧 Soluções para Problemas
- [Problema e solução]

## 📝 Observações Finais
[Qualquer observação adicional relevante]

IMPORTANTE:
- Seja conciso mas completo
- Priorize informações importantes
- Se houver conflitos (ex: "conflito no ambiente de testes"), sugira soluções práticas
- Use emojis para melhorar a legibilidade
- O relatório deve ser em português brasileiro
`
}

func (j *JarvisAgent) RunDaily() (string, error) {
	log.Println("Iniciando análise diária...")

	notesContent, err := j.reader.ReadTodayNotes()
	if err != nil {
		return "", fmt.Errorf("erro ao ler anotações: %w", err)
	}

	log.Printf("Anotações lidas: %d caracteres", len(notesContent))

	prompt := fmt.Sprintf(`
Analise as seguintes anotações diárias e gere um relatório estruturado:

%s

Gere o relatório completo em português brasileiro seguindo o formato definido nas suas instruções.
`, notesContent)

	userMsg := genai.NewContentFromText(prompt, "user")

	var report string
	for event, err := range j.runnerx.Run(context.Background(), "jarvis", "session-001", userMsg, agent.RunConfig{}) {
		if err != nil {
			return "", fmt.Errorf("erro ao executar agente: %w", err)
		}
		if event.Content != nil && len(event.Content.Parts) > 0 {
			if event.Content.Parts[0].Text != "" {
				report = event.Content.Parts[0].Text
			}
		}
	}

	if report == "" {
		return "", fmt.Errorf("nenhuma resposta do agente")
	}

	log.Printf("Relatório gerado: %d caracteres", len(report))

	if err := j.writer.WriteReport(report); err != nil {
		return "", fmt.Errorf("erro ao salvar relatório: %w", err)
	}

	log.Printf("Relatório salvo em: %s", j.writer.GetReportPath())

	return report, nil
}
