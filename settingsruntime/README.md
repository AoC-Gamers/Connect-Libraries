# settingsruntime

Biblioteca Go compartida para cargar y validar settings runtime desde Connect-Core durante el arranque de servicios.

Versión actual: `0.1.0`.

## Objetivo

- Esperar a que Connect-Core esté disponible (`/health`) con reintentos.
- Cargar claves requeridas desde `/core/internal/settings/{ENTITY}/{KEY}/value`.
- Validar formato/valor de cada setting antes de iniciar la aplicación.
- Entregar un mapa de valores listo para priming de caché local.

## Instalación

```bash
go get github.com/AoC-Gamers/connect-libraries/settingsruntime
```

## API principal

- `NewCoreHTTPClient(baseURL, apiKey string, timeout time.Duration) (*CoreHTTPClient, error)`
- `(*CoreHTTPClient).Health(ctx context.Context) error`
- `(*CoreHTTPClient).GetSettingValue(ctx context.Context, entity, key string) (string, error)`
- `Bootstrap(ctx context.Context, cfg BootstrapConfig) (map[string]string, error)`
- `CacheKey(entity, key string) string`
- `BoolValidator()`, `IntValidator(min,max)`, `NonEmptyValidator()`

## Ejemplo de uso

```go
client, err := settingsruntime.NewCoreHTTPClient(coreURL, apiKey, 5*time.Second)
if err != nil {
    return err
}

required := []settingsruntime.KeySpec{
    {Entity: "CONFIG", Key: "rate_limit.enabled", Validate: settingsruntime.BoolValidator()},
    {Entity: "CONFIG", Key: "rate_limit.rps", Validate: settingsruntime.IntValidator(1, 10000)},
}

values, err := settingsruntime.Bootstrap(ctx, settingsruntime.BootstrapConfig{
    ServiceName: "connect-auth",
    Retry: settingsruntime.RetryConfig{
        InitialInterval: time.Second,
        MaxInterval:     30 * time.Second,
        MaxElapsed:      2 * time.Minute,
        JitterFactor:    0.2,
    },
    Required:    required,
    HealthCheck: client.Health,
    Getter:      client.GetSettingValue,
})
if err != nil {
    return err
}

_ = values // priming de cache en el servicio
```

## Estructura interna

- `types.go`: tipos públicos (`KeySpec`, `RetryConfig`, `BootstrapConfig`).
- `client.go`: cliente HTTP de Core.
- `bootstrap.go`: orquestación de bootstrap y validación.
- `retry.go`: backoff/jitter y espera de readiness.
- `validators.go`: validadores reutilizables.
- `keys.go`: helpers de clave normalizada.
