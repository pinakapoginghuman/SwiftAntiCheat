package net.swiftsmp.anticheat.managers;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.entity.Player;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class BanManager {
    private final SwiftACPlugin plugin;
    private static final Pattern DURATION_PATTERN = Pattern.compile("(\\d+)([dhms]|permanent)", Pattern.CASE_INSENSITIVE);

    public BanManager(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    public void banPlayer(Player player, String durationStr, String reason) {
        Matcher m = DURATION_PATTERN.matcher(durationStr);
        if (!m.matches()) {
            Bukkit.dispatchCommand(Bukkit.getConsoleSender(),
                "ban " + player.getName() + " " + reason);
            return;
        }

        String numStr = m.group(1);
        String unit = m.group(2).toLowerCase();
        int amount = Integer.parseInt(numStr);

        if (unit.equals("permanent")) {
            Bukkit.dispatchCommand(Bukkit.getConsoleSender(),
                "ban " + player.getName() + " " + reason);
            return;
        }

        ChronoUnit chronoUnit;
        switch (unit) {
            case "d": chronoUnit = ChronoUnit.DAYS; break;
            case "h": chronoUnit = ChronoUnit.HOURS; break;
            case "m": chronoUnit = ChronoUnit.MINUTES; break;
            case "s": chronoUnit = ChronoUnit.SECONDS; break;
            default: chronoUnit = ChronoUnit.DAYS;
        }

        long expires = Instant.now().plus(amount, chronoUnit).toEpochMilli();
        Bukkit.dispatchCommand(Bukkit.getConsoleSender(),
            "ban " + player.getName() + " " + expires + " " + reason);
    }
}
