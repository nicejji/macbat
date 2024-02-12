use std::{io::BufRead, process::Command};

const RESET: &str = "\x1b[0m";
const BOLD: &str = "\x1b[1m";
const NORMAL: &str = "\x1b[22m";
const ITALIC: &str = "\x1b[3m";
const BLUE: &str = "\x1b[34m";
const RED: &str = "\x1b[31m";
const YELLOW: &str = "\x1b[33m";
const CYAN: &str = "\x1b[36m";
const GRAY: &str = "\x1b[90m";

#[derive(Default)]
struct Info {
    temperature: f64,
    cycles_count: isize,
    design_capacity: isize,
    max_capacity: isize,
    current_capacity: isize,
    is_charging: bool,
}

impl Info {
    pub fn print(&self) {
        let health_percent = self.max_capacity as f64 / self.design_capacity as f64 * 100.0;
        let charge_percent = self.current_capacity as f64 / self.max_capacity as f64 * 100.0;
        let charging_symbol = if self.is_charging { '⇡' } else { '⇣' };
        println!(
            "{}",
            format_opt(
                "Raw Health",
                &format!(
                    "{:.2}% ({}/{} mAh)",
                    health_percent, self.max_capacity, self.design_capacity
                ),
                "Raw battery health",
                BLUE
            )
        );
        println!(
            "{}",
            format_opt(
                "Temperature",
                &format!("{:.2} °C", self.temperature),
                "Battery temperature",
                RED
            )
        );
        println!(
            "{}",
            format_opt(
                "Cycles count",
                &format!("{}", self.cycles_count),
                "Cycles count",
                YELLOW
            )
        );
        println!(
            "{}",
            format_opt(
                "Charge info",
                &format!(
                    "{:.2} ({}/{} mAh) {}",
                    charge_percent, self.current_capacity, self.max_capacity, charging_symbol
                ),
                "Charge percent",
                CYAN
            )
        );
    }
}

fn format_opt(name: &str, value: &str, desc: &str, color: &str) -> String {
    return format!("{color}{NORMAL}{name}\t{BOLD}{value}  {GRAY}{NORMAL}{ITALIC}{desc}{RESET}");
}

fn main() {
    let output = Command::new("ioreg")
        .args(["-w", "0", "-r", "-c", "AppleSmartBattery"])
        .output()
        .expect("Failed to run ioreg");
    let mut info = Info::default();
    for line in output.stdout.lines().map(|l| l.unwrap()) {
        let parts: Vec<&str> = line.split("=").collect();
        let name = parts.get(0);
        let value = parts.get(1);
        match name {
            Some(name) => match name.trim().trim_matches('"') {
                "IsCharging" => info.is_charging = value.unwrap().trim() == "Yes",
                "AppleRawMaxCapacity" => info.max_capacity = value.unwrap().trim().parse().unwrap(),
                "AppleRawCurrentCapacity" => {
                    info.current_capacity = value.unwrap().trim().parse().unwrap();
                }

                "DesignCapacity" => info.design_capacity = value.unwrap().trim().parse().unwrap(),
                "CycleCount" => info.cycles_count = value.unwrap().trim().parse().unwrap(),
                "Temperature" => {
                    info.temperature = value.unwrap().trim().parse::<f64>().unwrap() / 100.0
                }
                _ => {}
            },
            None => {}
        }
    }
    info.print();
}
