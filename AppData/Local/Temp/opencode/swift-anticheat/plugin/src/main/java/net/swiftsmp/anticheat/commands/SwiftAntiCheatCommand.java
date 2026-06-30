package net.swiftsmp.anticheat.commands;

import net.swiftsmp.anticheat.SwiftACPlugin;
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

        String scanUrl = plugin.getConfig().getString("scan.base-url", "https://swift-anti-cheat.vercel.app")
            + "/scan";

        if (!plugin.getScreenshareManager().isFrozen(target.getUniqueId())) {
            plugin.getScreenshareManager().freezePlayer(target, staff);
        }

        for (String line : plugin.getConfig().getStringList("scan.link-message")) {
            target.sendMessage(ChatColor.translateAlternateColorCodes('&',
                line.replace("{link}", scanUrl)
                    .replace("{player}", target.getName())
                    .replace("{staff}", staff.getName())));
        }

        staff.sendMessage(ChatColor.GREEN + "Scan link sent to " + target.getName() + "!");
        staff.sendMessage(ChatColor.GRAY + "Player will get a report code after scanning.");
        staff.sendMessage(ChatColor.GRAY + "Go to " + ChatColor.AQUA + scanUrl + ChatColor.GRAY + " to enter the code and view results.");

        return true;
    }
}
