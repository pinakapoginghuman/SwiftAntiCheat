package net.swiftsmp.anticheat.commands;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

public class ScreenshareCommand implements CommandExecutor {
    private final SwiftACPlugin plugin;

    public ScreenshareCommand(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @Override
    public boolean onCommand(CommandSender sender, Command cmd, String label, String[] args) {
        if (!(sender instanceof Player staff)) {
            sender.sendMessage("Only players can use this command.");
            return true;
        }

        if (args.length < 1) {
            staff.sendMessage(ChatColor.RED + "Usage: /screenshare <player>");
            return true;
        }

        Player target = Bukkit.getPlayer(args[0]);
        if (target == null) {
            staff.sendMessage(ChatColor.RED + "Player not found.");
            return true;
        }

        if (target.equals(staff)) {
            staff.sendMessage(ChatColor.RED + "You cannot screenshare yourself.");
            return true;
        }

        if (plugin.getScreenshareManager().isFrozen(target.getUniqueId())) {
            staff.sendMessage(ChatColor.RED + "That player is already being screenshared.");
            return true;
        }

        plugin.getScreenshareManager().freezePlayer(target, staff);
        return true;
    }
}
