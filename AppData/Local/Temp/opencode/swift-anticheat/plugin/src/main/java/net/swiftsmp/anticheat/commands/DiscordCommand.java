package net.swiftsmp.anticheat.commands;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

public class DiscordCommand implements CommandExecutor {
    private final SwiftACPlugin plugin;

    public DiscordCommand(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String label, String[] args) {
        if (!(sender instanceof Player)) {
            sender.sendMessage("Only players can use this command.");
            return true;
        }

        String discord = plugin.getConfig().getString("discord.invite", "https://discord.gg/swiftsmp");
        for (String line : plugin.getConfig().getStringList("discord.message")) {
            sender.sendMessage(ChatColor.translateAlternateColorCodes('&',
                line.replace("{discord}", discord)));
        }

        return true;
    }
}
