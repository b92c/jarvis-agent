# Jarvis - Assistente de Team/Tech Lead

Assistente pessoal automatizado que analiza suas anotações diárias e gera relatórios estruturados.

## Funcionalidades

- **Execução automática**: Roda todos os dias às 12:00 e 18:00
- **Leitura de anotações**: Lê o arquivo `today.md` do diretório de notas
- **Análise com IA**: Usa Google Gemini (gemini-3-pro-preview) para analisar as anotações
- **Geração de relatório**: Cria `jarvis-report.md` com:
  - Principais pontos do dia
  - Pontos de atenção
  - Ações práticas sugeridas
  - Soluções para problemas
- **Envio por email**: Envia o relatório por email automaticamente

## Configuração

1. Clone o repositório e configure as variáveis de ambiente no arquivo `.env`:

```bash
# API Key do Google AI Studio (Gemini)
GOOGLE_API_KEY=sua_chave_aqui

# Paths (opcional - padrões já configurados)
TODAY_NOTES_PATH=/home/user/Documents/personal-notes/Trabalho/Gestão/Notes/today.md
REPORT_OUTPUT_PATH=/home/user/Documents/personal-notes/Trabalho/Gestão/Notes/jarvis-report.md

# Email (Gmail com App Password)
SMTP_PASS=sua_app_password_aqui
```

### Obtendo a Google API Key
1. Acesse [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Crie uma nova API key
3. Cole no arquivo `.env`

### Configurando App Password do Gmail
1. Ative a verificação em duas etapas na sua conta Google
2. Acesse [App Passwords](https://myaccount.google.com/apppasswords)
3. Crie uma nova senha para o aplicativo
4. Use essa senha no campo `SMTP_PASS`

## Executando

### Localmente
```bash
go run .
```

### Com Docker
```bash
docker-compose up -d
```

## Estrutura do Projeto

```
jarvis-agent/
├── main.go           # Entry point + scheduler
├── agent/
│   └── jarvis.go    # Definição do agente
├── tools/
│   ├── file_reader.go
│   └── file_writer.go
├── services/
│   └── email.go     # Serviço de email
├── config/
│   └── config.go    # Configurações
├── Dockerfile
├── docker-compose.yml
└── .env
```

## Formato do Relatório

O agente gera um relatório em Markdown com o seguinte formato:

```markdown
# 📋 Jarvis Report - [DATA]

## 🎯 Principais Pontos do Dia
- ...

## ⚠️ Pontos de Atenção
- ...

## 💡 Ações Práticas Sugeridas
- ...

## 🔧 Soluções para Problemas
- ...

## 📝 Observações Finais
- ...
```

## Licença

MIT
