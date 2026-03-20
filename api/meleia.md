**Cada linguagem precisa pagar o custo dela**

Se você adicionar uma linguagem, se pergunte:

    ela resolve um problema que as outras não resolvem?
    ou só tá deixando o sistema mais complexo?

**fluxo completo:**
```
GitHub faz POST /webhook
    → Gin recebe a requisição
        → WebhookHandler processa
            → Responde JSON pro GitHub


Por que essa estrutura é "limpa"?

Pasta                   Responsabilidade                Analogia POO

pkg/config              Sabe o que configurar           Classe de configuração
pkg/middleware          Sabe como proteger              Interceptors/Filters
internal/handler        Sabe o que responder            Controllers
cmd/api/main.go         Só monta as peças               Main/Application

Cada camada não sabe da existência das outras — o middleware não importa o handler, o handler não importa o config. Isso é o princípio de responsabilidade única (SRP) na prática.



**Resumo do que foi construído até agora:**
```
✅ API Go com Gin rodando
✅ CORS configurado para Cloudflare  
✅ Headers de segurança HTTP
✅ Validação HMAC-SHA256 dos webhooks do GitHub
✅ Config centralizada e limpa