package net.swiftsmp.anticheat.listeners;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.entity.Player;
import org.bukkit.event.EventHandler;
import org.bukkit.event.EventPriority;
import org.bukkit.event.Listener;
import org.bukkit.event.player.AsyncPlayerChatEvent;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.regex.Pattern;

public class ChatListener implements Listener {
    private final SwiftACPlugin plugin;
    private final Pattern CODE_PATTERN = Pattern.compile("^SWIFT-[A-Z0-9]{4}-[A-Z0-9]{4}$");
    private final HttpClient client;

    public ChatListener(SwiftACPlugin plugin) {
        this.plugin = plugin;
        this.client = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(10))
            .build();
    }

    @EventHandler(priority = EventPriority.HIGHEST)
    public void onPlayerChat(AsyncPlayerChatEvent event) {
        Player player = event.getPlayer();

        if (plugin.getScreenshareManager().isFrozen(player.getUniqueId())) {
            String message = event.getMessage().trim().toUpperCase();

            if (CODE_PATTERN.matcher(message).matches()) {
                event.setCancelled(false);
                event.setFormat(ChatColor.GRAY + "[Code] " + ChatColor.WHITE + "%2$s");
                event.getRecipients().clear();
                event.getRecipients().addAll(Bukkit.getOnlinePlayers());

                var session = plugin.getScreenshareManager().getSession(player.getUniqueId());
                String staffName = session != null ? Bukkit.getPlayer(session.getStaffUUID()) != null ? Bukkit.getPlayer(session.getStaffUUID()).getName() : null : null;
                String apiUrl = plugin.getConfig().getString("api.base-url", "https://swiftac-api.onrender.com");
                String dashUrl = plugin.getConfig().getString("scan.base-url", "https://swift-anti-cheat.vercel.app");

                String lookupUrl = dashUrl + "/scan?code=" + message;

                if (staffName != null) {
                    Player staff = Bukkit.getPlayer(staffName);
                    if (staff != null && staff.isOnline()) {
                        staff.sendMessage("");
                        staff.sendMessage(ChatColor.GREEN + "═══════════════════════════════════");
                        staff.sendMessage(ChatColor.GREEN + "  " + player.getName() + " submitted their scan code!");
                        staff.sendMessage(ChatColor.GREEN + "  Code: " + ChatColor.AQUA + message);
                        staff.sendMessage(ChatColor.GREEN + "  " + ChatColor.UNDERLINE + lookupUrl);
                        staff.sendMessage(ChatColor.GREEN + "═══════════════════════════════════");
                        staff.sendMessage("");
                    }
                }

                String finalMessage = message;
                Bukkit.getScheduler().runTaskAsynchronously(plugin, () -> {
                    checkReportCode(apiUrl, finalMessage, player, staffName);
                });
            } else {
                event.setCancelled(true);
                player.sendMessage(ChatColor.RED + "Chat is disabled while frozen.");
                player.sendMessage(ChatColor.GRAY + "Type your report code (e.g. " + ChatColor.AQUA + "SWIFT-XXXX-XXXX" + ChatColor.GRAY + ") to submit it.");
            }
            return;
        }

        event.getRecipients().removeIf(r ->
            plugin.getScreenshareManager().isFrozen(r.getUniqueId()));
    }

    private void checkReportCode(String apiUrl, String code, Player player, String staffName) {
        try {
            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(apiUrl + "/api/reports/" + code))
                .timeout(Duration.ofSeconds(15))
                .GET()
                .build();

            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

            Bukkit.getScheduler().runTask(plugin, () -> {
                if (response.statusCode() == 200) {
                    player.sendMessage(ChatColor.GREEN + "Your scan results have been submitted! Staff will review them.");
                    plugin.getScreenshareManager().unfreezePlayer(player);

                    if (staffName != null) {
                        Player staff = Bukkit.getPlayer(staffName);
                        if (staff != null && staff.isOnline()) {
                            staff.sendMessage(ChatColor.GREEN + player.getName() + "'s scan code verified! Results are ready.");
                        }
                    }
                } else {
                    player.sendMessage(ChatColor.RED + "Report code not found. Make sure you typed it correctly.");
                }
            });
        } catch (Exception e) {
            player.sendMessage(ChatColor.RED + "Could not verify report code. The API may be waking up.");
            player.sendMessage(ChatColor.GRAY + "Try again in a few seconds, or tell staff the code directly.");
        }
    }
}
