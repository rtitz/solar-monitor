# Solar-Monitor

### Build the binaries & run the Ansible playbook:
```zsh
cd src && bash build.sh && cd .. && ansible-playbook -i inventory main.yml
```

### Example
```zsh
journalctl -u solar-monitor -f
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
