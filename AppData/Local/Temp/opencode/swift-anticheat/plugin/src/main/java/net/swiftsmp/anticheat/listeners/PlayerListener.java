package net.swiftsmp.anticheat.listeners;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.entity.Player;
import org.bukkit.event.EventHandler;
import org.bukkit.event.EventPriority;
import org.bukkit.event.Listener;
import org.bukkit.event.player.PlayerJoinEvent;
import org.bukkit.event.player.PlayerMoveEvent;
import org.bukkit.event.player.PlayerQuitEvent;
import org.bukkit.event.player.PlayerTeleportEvent;

public class PlayerListener implements Listener {
    private final SwiftACPlugin plugin;

    public PlayerListener(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    @EventHandler(priority = EventPriority.LOWEST)
    public void onPlayerMove(PlayerMoveEvent event) {
        Player player = event.getPlayer();
        if (plugin.getScreenshareManager().isFrozen(player.getUniqueId())) {
            var session = plugin.getScreenshareManager().getSession(player.getUniqueId());
            if (session != null) {
                if (event.getTo().getX() != session.getFrozenLocation().getX()
                    || event.getTo().getZ() != session.getFrozenLocation().getZ()) {
                    event.setTo(session.getFrozenLocation().clone());
                }
            }
        }
    }

    @EventHandler(priority = EventPriority.LOWEST)
    public void onPlayerTeleport(PlayerTeleportEvent event) {
        Player player = event.getPlayer();
        if (plugin.getScreenshareManager().isFrozen(player.getUniqueId())) {
            var session = plugin.getScreenshareManager().getSession(player.getUniqueId());
            if (session != null && event.getCause() != PlayerTeleportEvent.TeleportCause.UNKNOWN) {
                event.setCancelled(true);
            }
        }
    }

    @EventHandler(priority = EventPriority.MONITOR)
    public void onPlayerQuit(PlayerQuitEvent event) {
        Player player = event.getPlayer();
        if (plugin.getScreenshareManager().isFrozen(player.getUniqueId())) {
            plugin.getScreenshareManager().handleDodge(player);
        }
    }
}
