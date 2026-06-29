package net.swiftsmp.anticheat.commands;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

public class UnscreenshareCommand implements CommandExecutor {
    private final SwiftACPlugin plugin;

    public UnscreenshareCommand(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String label, String[] args) {
        if (!(sender instanceof Player)) {
            sender.sendMessage("Only players can use this command.");
            return true;
        }

        if (args.length < 1) {
            sender.sendMessage(ChatColor.RED + "Usage: /unscreenshare <player>");
            return true;
        }

        Player target = Bukkit.getPlayer(args[0]);
        if (target == null) {
            sender.sendMessage(ChatColor.RED + "Player not found.");
            return true;
        }

        if (!plugin.getScreenshareManager().isFrozen(target.getUniqueId())) {
            sender.sendMessage(ChatColor.RED + "That player is not frozen.");
            return true;
        }

        plugin.getScreenshareManager().unfreezePlayer(target);
        sender.sendMessage(ChatColor.GREEN + "Unfrozen " + target.getName() + ".");
        return true;
    }
}
