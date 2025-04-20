# Solar-Monitor


### Build the binaries & run the Ansible playbook:
```zsh
cd src && bash build.sh && cd .. && ansible-playbook -i inventory main.yml
```

### Example
```zsh
journalctl -u solar-monitor -f
```