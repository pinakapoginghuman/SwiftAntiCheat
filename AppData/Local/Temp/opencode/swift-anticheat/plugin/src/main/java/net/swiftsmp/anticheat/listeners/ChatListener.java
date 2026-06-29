package net.swiftsmp.anticheat.listeners;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.entity.Player;
import org.bukkit.event.EventHandler;
import org.bukkit.event.EventPriority;
import org.bukkit.event.Listener;
import org.bukkit.event.player.AsyncPlayerChatEvent;

public class ChatListener implements Listener {
    private final SwiftACPlugin plugin;

    public ChatListener(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @EventHandler(priority = EventPriority.HIGHEST)
    public void onPlayerChat(AsyncPlayerChatEvent event) {
        Player player = event.getPlayer();

        if (plugin.getScreenshareManager().isFrozen(player.getUniqueId())) {
            event.setCancelled(true);
            player.sendMessage(org.bukkit.ChatColor.RED + "Chat is disabled while being screenshared.");
            return;
        }

        event.getRecipients().removeIf(r ->
            plugin.getScreenshareManager().isFrozen(r.getUniqueId()));
    }
}
