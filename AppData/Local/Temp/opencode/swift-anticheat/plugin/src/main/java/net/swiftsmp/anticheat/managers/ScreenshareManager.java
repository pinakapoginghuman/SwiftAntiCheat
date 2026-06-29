package net.swiftsmp.anticheat.managers;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.Location;
import org.bukkit.entity.Player;
import org.bukkit.scheduler.BukkitTask;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

public class ScreenshareManager {
    private final SwiftACPlugin plugin;
    private final Map<UUID, ScreenshareSession> sessions = new ConcurrentHashMap<>();
    private final Map<UUID, BukkitTask> freezeTasks = new ConcurrentHashMap<>();

    public ScreenshareManager(SwiftACPlugin plugin) {
        this.plugin = plugin;
    }

    public void freezePlayer(Player target, Player staff) {
        UUID uuid = target.getUniqueId();
        if (sessions.containsKey(uuid)) return;

        boolean freezeAi = plugin.getConfig().getBoolean("screenshare.freeze-ai", true);
        Location frozenLocation = target.getLocation().clone();

        sessions.put(uuid, new ScreenshareSession(uuid, staff.getUniqueId(), frozenLocation));

        if (freezeAi) {
            target.setWalkSpeed(0f);
            target.setFlySpeed(0f);
            target.setAllowFlight(true);
            target.setFlying(true);
            target.setCollidable(false);
            target.setInvulnerable(true);

            BukkitTask task = Bukkit.getScheduler().runTaskTimer(plugin, () -> {
                Player p = Bukkit.getPlayer(uuid);
                if (p != null && sessions.containsKey(uuid)) {
                    p.teleport(frozenLocation);
                    p.setWalkSpeed(0f);
                    p.setFlySpeed(0f);
                    p.setAllowFlight(true);
                    p.setFlying(true);
                }
            }, 0L, 5L);
            freezeTasks.put(uuid, task);
        }

        List<String> msg = plugin.getConfig().getStringList("screenshare.freeze-message");
        for (String line : msg) {
            target.sendMessage(org.bukkit.ChatColor.translateAlternateColorCodes('&',
                line.replace("{staff}", staff.getName())
                    .replace("{duration}", plugin.getConfig().getString("screenshare.dodge-ban-duration", "7d"))));
        }
        staff.sendMessage(org.bukkit.ChatColor.translateAlternateColorCodes('&',
            plugin.getConfig().getString("screenshare.staff-freeze-notify", "&a{player} has been frozen.")
                .replace("{player}", target.getName())));
    }

    public void unfreezePlayer(Player target) {
        UUID uuid = target.getUniqueId();
        if (!sessions.containsKey(uuid)) return;

        sessions.remove(uuid);
        BukkitTask task = freezeTasks.remove(uuid);
        if (task != null) task.cancel();

        target.setWalkSpeed(0.2f);
        target.setFlySpeed(0.1f);
        target.setAllowFlight(false);
        target.setFlying(false);
        target.setCollidable(true);
        target.setInvulnerable(false);
        target.setFallDistance(0f);

        target.sendMessage(org.bukkit.ChatColor.translateAlternateColorCodes('&',
            plugin.getConfig().getString("screenshare.unfreeze-message", "&aUnfrozen.")));
    }

    public void unfreezeAll() {
        for (UUID uuid : new HashSet<>(sessions.keySet())) {
            Player p = Bukkit.getPlayer(uuid);
            if (p != null) unfreezePlayer(p);
        }
        for (BukkitTask task : freezeTasks.values()) task.cancel();
        freezeTasks.clear();
        sessions.clear();
    }

    public boolean isFrozen(UUID uuid) {
        return sessions.containsKey(uuid);
    }

    public ScreenshareSession getSession(UUID uuid) {
        return sessions.get(uuid);
    }

    public void handleDodge(Player player) {
        UUID uuid = player.getUniqueId();
        ScreenshareSession session = sessions.get(uuid);
        if (session == null) return;

        unfreezePlayer(player);
        String duration = plugin.getConfig().getString("screenshare.dodge-ban-duration", "7d");
        String reason = plugin.getConfig().getString("screenshare.dodge-ban-reason", "Dodging screenshare");
        Bukkit.getScheduler().runTask(plugin, () -> plugin.getBanManager().banPlayer(player, duration, reason));

        Bukkit.broadcast(org.bukkit.ChatColor.translateAlternateColorCodes('&',
            plugin.getConfig().getString("screenshare.dodge-notify", "&c{player} dodged!")
                .replace("{player}", player.getName())
                .replace("{duration}", duration)), "swiftac.staff");
    }

    public static class ScreenshareSession {
        private final UUID playerUUID;
        private final UUID staffUUID;
        private final Location frozenLocation;
        private final long startTime;

        public ScreenshareSession(UUID playerUUID, UUID staffUUID, Location frozenLocation) {
            this.playerUUID = playerUUID;
            this.staffUUID = staffUUID;
            this.frozenLocation = frozenLocation;
            this.startTime = System.currentTimeMillis();
        }

        public UUID getPlayerUUID() { return playerUUID; }
        public UUID getStaffUUID() { return staffUUID; }
        public Location getFrozenLocation() { return frozenLocation; }
        public long getStartTime() { return startTime; }
    }
}
