# 🚗 Mi Coche Ideal - Backend Microservice (Golang)

Microservicio REST en Go para gestionar el inventario y la venta de vehículos usados para un concesionario ficticio.

---

## 🧱 Tecnologías

- **Lenguaje:** Go 1.21
- **Framework:** Gin
- **ORM:** GORM + PostgreSQL
- **Autenticación:** JWT + Refresh Tokens
- **Caché (opcional):** Redis
- **Observabilidad:** Prometheus + Zerolog
- **Contenedores:** Docker + Docker Compose
- **Pruebas:** Go testing + Testify

---

## ▶️ Cómo correr el proyecto

### 1. Clonar el repositorio

```bash
git clone https://github.com/JesusOsorioJ/mi-coche-ideal.git
cd mi-coche-ideal
```

### 2. Crear archivo `.env`

Puedes crear un archivo `.env` con tus variables personalizadas o copiar el archivo de ejemplo incluido:

```bash
cp .env.example .env
```

Ejemplo de contenido:

```env
PORT=3000
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mi_coche_ideal
JWT_SECRET=supersecretokey
```

### 3. Levantar con Docker Compose

```bash
docker-compose up --build
```

---

## ✅ Endpoints disponibles

### 🔐 Autenticación
- `POST /auth/signup` – Registro
- `POST /auth/login` – Login (devuelve `access_token` y `refresh_token`)

### 🚗 Vehículos (requiere token)
- `POST /api/vehicles`
- `GET /api/vehicles` – con filtros (`?brand=Audi&price_min=10000`)
- `PUT /api/vehicles/:id`
- `DELETE /api/vehicles/:id`

### 🧾 Órdenes de venta (requiere token)
- `POST /api/orders`
- `GET /api/orders`
- `GET /api/orders/:id`
- `PUT /api/orders/:id/status` – Actualiza estado (`pagada`, `entregada`)

### 📊 Observabilidad
- `GET /healthz` – Health check
- `GET /metrics` – Métricas Prometheus

---

## 🧪 Pruebas

### Pruebas automatizadas

```bash
go test ./... -cover
```

o con reporte visual:

```bash
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

📈 **Objetivo de cobertura:** mayor al **70%**  
Las pruebas cubren los módulos de **autenticación**, **gestión de vehículos** y **órdenes**.

> ⚠️ **Importante:**  
> Si levanta la base de datos con Docker, asegúrate de ajustar la variable de entorno del **host de la base de datos** a `localhost` **antes de ejecutar los tests**.

### Pruebas manuales

Incluye un archivo `test_api.rest` compatible con [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) en VSCode para probar fácilmente los endpoints sin Postman.

---

## ⏱️ Concurrencia

Se ejecuta cada minuto una rutina que:

- Lee precios desde `price_updates.csv`
- Aplica actualizaciones concurrentes usando **goroutines** y **mutex**
- Utiliza **transacciones** para mantener la consistencia

---

## 🧠 Decisiones técnicas

- Uso de GORM + AutoMigrate para agilidad durante el desarrollo
- Patrón handler-service-repository para separación de responsabilidades
- Logs en formato JSON con `zerolog`
- Uso de `robfig/cron` para agendado confiable de tareas
- JWT y Refresh Tokens separados para mejor control de sesiones
- Rutas `/metrics` y `/healthz` accesibles sin autenticación

---

## 📂 Archivos útiles

- `.env.example`: plantilla base para crear tu archivo `.env`
- `test_api.rest`: colección de endpoints para pruebas manuales vía REST Client

---

## 📝 Post-mortem

- Con más tiempo, integraría **Redis** para caché y mutex distribuidos.
- Añadiría pruebas de estrés (load testing) y CI/CD con GitHub Actions.
- Estandarizaría el manejo de errores con un paquete común.
- Implementaría rollback automático si una rutina de actualización falla parcialmente.

---
