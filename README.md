# Solar-Monitor

### Build the binaries & run the Ansible playbook:
```zsh
cd src && bash build.sh && cd .. && ansible-playbook -i inventory main.yml
# If podman container:
ansible-playbook -i inventory main.yml
```

### Example
```zsh
journalctl -u solar-monitor -f
```

### Build Container
```zsh
podman build -t solar-monitor .
podman run -d -p 0.0.0.0:2112:2112 solar-monitor
```

## Prometheus Integration

The application exposes metrics for Prometheus at `http://localhost:2112/metrics`.

### Metrics Available:
- `solar_gen_a_power_watts` - Solar generator A power
- `solar_gen_b_power_watts` - Solar generator B power  
- `inverter_ac_power_watts` - Inverter AC power
- `total_solar_power_watts` - Total solar power
- `total_house_consumption_watts` - Total house consumption
- `total_grid_power_watts` - Grid power (+ consuming, - feeding)
- `using_solar_power_only` - Using solar only (1=true, 0=false)
- `solar_power_available` - Solar available (1=true, 0=false)

### Prometheus Configuration:
```yaml
scrape_configs:
  - job_name: 'solar-monitor'
    static_configs:
      - targets: ['localhost:2112']
    scrape_interval: 5s
```

## Build container

```zsh
mkdir -p /data/solar-monitor ; chmod 1777 /data/solar-monitor
TMPDIR=/tmp docker build -t solar-monitor:latest .
TMPDIR=/tmp docker run -d --name solar-monitor -p 2112:2112 -v /data/solar-monitor:/home/appuser/solar-monitor:Z --restart unless-stopped solar-monitor:latest

# OR run in foreground
TMPDIR=/tmp docker run -p 2112:2112 -v /data/solar-monitor:/home/appuser/solar-monitor:Z --rm solar-monitor:latest
```

---

## ⚙️ Configuration Parameter Reference

The metrics monitoring engine is completely dynamic and reads its environment configurations on startup. If an environment variable is missing, it safely drops back to its built-in out-of-the-box defaults.

### Core Service Configuration Variables

| Environment Variable | Description | Expected Type | Default Value |
| :--- | :--- | :--- | :--- |
| **`INVERTER_IP`** | Network destination IP address of your physical RCT power inverter. (Port is fixed in RCT-API) | String | `192.168.1.17` |
| **`PROMETHEUS_PORT`** | Active network port binding configuration for exposing Prometheus metrics. | String | `2112` |
| **`FORCE_IPV4`** | Binds network socket connection layer listening strictly to the IPv4 stack. | Boolean | `true` |
| **`FETCH_INTERVAL_MS`** | Direct hardware register polling execution rate interval in milliseconds. | Integer | `1000` |
| **`RETRY_INTERVAL_MS`** | Network error collection buffer delay before initiating an API query retry. | Integer | `250` |
| **`MAX_RETRIES`** | Maximum sequential hardware read failures before marking active values as invalid. | Integer | `60` |
| **`START_DELAY_MS`** | Initial runtime pause before bootstrapping background network threads. | Integer | `2000` |
| **`OUTPUT_REFRESH_INTERVAL_MS`** | Core engine output state refresh cadence (and interval for flushing tracking JSON files). | Integer | `2000` |
| **`HOME_ASSISTANT_UPDATE_INTERVAL_ITERATIONS`** | Multiplier factor for sending states (`Iterations` × `OUTPUT_REFRESH_INTERVAL_MS`, e.g., 5 × 2000ms = 10s). | Integer | `5` |
| **`HOME_ASSISTANT_LOG_API_CALL_RESULT_200`** | Toggles console telemetry write logs for successful Home Assistant API push executions. | Boolean | `true` |

---

### Home Assistant Integration Configuration Variables

| Environment Variable | Description | Expected Type | Default Value |
| :--- | :--- | :--- | :--- |
| **`URL_HOME_ASSISTANT`** | Target API URL endpoint of your active Home Assistant server instance. | String | `http://192.168.1.8:8123` |
| **`TOKEN_HOME_ASSISTANT`** | Long-Lived Access Token generated from Home Assistant for secure REST authentication. | String | *(Uses dynamic secure token)* |
| **`API_BOOL_SOLAR_POWER_ONLY`** | Target entity API path mapping for pure solar calculation tracking tracking. | String | `/api/states/input_boolean.solar_power_only` |
| **`API_BOOL_SOLAR_POWER_AVAILABLE`** | Target entity API path mapping for solar generation availability state indicators. | String | `/api/states/input_boolean.solar_power_available` |
| **`API_INT_SOLAR_GEN_A_POWER_W`** | Target entity API path mapping for physical Solar String Generator A telemetry input. | String | `/api/states/input_number.solar_gen_a_power_w` |
| **`API_INT_SOLAR_GEN_B_POWER_W`** | Target entity API path mapping for physical Solar String Generator B telemetry input. | String | `/api/states/input_number.solar_gen_b_power_w` |
| **`API_INT_SOLAR_GEN_TOTAL_POWER_W`** | Target entity API path mapping for collective Inverter PV production tracking metrics. | String | `/api/states/input_number.solar_gen_total_power_w` |
| **`API_INT_PV_DELIVERY`** | Target entity API path mapping for total real-time solar generation delivery numbers. | String | `/api/states/input_number.pv_erzeugung` |
| **`API_INT_HOUSE_CONSUMPTION`** | Target entity API path mapping for total calculated residential load energy consumption metrics. | String | `/api/states/input_number.stromverbrauch` |
| **`API_INT_GRID_FEED`** | Target entity API path mapping for real-time grid energy balancing metrics (Feed-in/Import). | String | `/api/states/input_number.pv_einspeisung` |
