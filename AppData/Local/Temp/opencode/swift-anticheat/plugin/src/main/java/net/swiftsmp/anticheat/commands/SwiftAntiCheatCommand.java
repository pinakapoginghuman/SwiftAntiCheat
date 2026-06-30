package net.swiftsmp.anticheat.commands;

import net.swiftsmp.anticheat.SwiftACPlugin;
import net.swiftsmp.anticheat.managers.APIManager;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

public class SwiftAntiCheatCommand implements CommandExecutor {
    private final SwiftACPlugin plugin;

    public SwiftAntiCheatCommand(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String label, String[] args) {
        if (!(sender instanceof Player staff)) {
            sender.sendMessage("Only players can use this command.");
            return true;
        }

        if (args.length < 1) {
            staff.sendMessage(ChatColor.RED + "Usage: /swiftanticheat <player>");
            return true;
        }

        Player target = Bukkit.getPlayer(args[0]);
        if (target == null) {
            staff.sendMessage(ChatColor.RED + "Player not found.");
            return true;
        }

        staff.sendMessage(ChatColor.GRAY + "Generating scan link for " + target.getName() + "...");

        Bukkit.getScheduler().runTaskAsynchronously(plugin, () -> {
            APIManager.ScanResponse scanResponse = plugin.getAPIManager().createScan(
                target.getName(),
                target.getUniqueId().toString(),
                staff.getName()
            );

            Bukkit.getScheduler().runTask(plugin, () -> {
                if (scanResponse == null) {
                    staff.sendMessage(ChatColor.RED + "Failed to generate scan link. API is unreachable.");
                    return;
                }

                String scanUrl = plugin.getConfig().getString("scan.base-url", "https://swiftac-scan.vercel.app")
                    + "/scan?id=" + scanResponse.getId();

                if (!plugin.getScreenshareManager().isFrozen(target.getUniqueId())) {
                    plugin.getScreenshareManager().freezePlayer(target, staff);
                }

                for (String line : plugin.getConfig().getStringList("scan.link-message")) {
                    target.sendMessage(ChatColor.translateAlternateColorCodes('&',
                        line.replace("{link}", scanUrl)));
                }

                staff.sendMessage(ChatColor.GREEN + "Scan link generated!");
                staff.sendMessage(ChatColor.GRAY + "Link: " + ChatColor.AQUA + scanUrl);
            });
        });

        return true;
    }
}
