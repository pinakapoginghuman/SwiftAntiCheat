package net.swiftsmp.anticheat.commands;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;
import org.bukkit.ChatColor;
import org.bukkit.command.Command;
import org.bukkit.command.CommandExecutor;
import org.bukkit.command.CommandSender;
import org.bukkit.entity.Player;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

public class SwiftAntiCheatCommand implements CommandExecutor {
    private final SwiftACPlugin plugin;
    private final HttpClient httpClient;
    private final Gson gson = new Gson();

    public SwiftAntiCheatCommand(SwiftACPlugin plugin) {
        this.plugin = plugin;
        this.httpClient = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(10))
            .build();
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
            try {
                String baseUrl = plugin.getConfig().getString("api.base-url", "http://localhost:3000");
                String apiKey = plugin.getConfig().getString("api.api-key", "");

                JsonObject body = new JsonObject();
                body.addProperty("playerName", target.getName());
                body.addProperty("playerUUID", target.getUniqueId().toString());
                body.addProperty("staffName", staff.getName());

                HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(baseUrl + "/api/scans/create"))
                    .header("Content-Type", "application/json")
                    .header("Authorization", "Bearer " + apiKey)
                    .POST(HttpRequest.BodyPublishers.ofString(gson.toJson(body)))
                    .timeout(Duration.ofSeconds(10))
                    .build();

                HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

                if (response.statusCode() == 200 || response.statusCode() == 201) {
                    JsonObject json = gson.fromJson(response.body(), JsonObject.class);
                    String scanId = json.get("id").getAsString();
                    String scanUrl = plugin.getConfig().getString("scan.base-url", "https://swiftac-scan.vercel.app")
                        + "/scan/" + scanId;

                    Bukkit.getScheduler().runTask(plugin, () -> {
                        staff.sendMessage(ChatColor.GREEN + "Scan link generated for " + target.getName() + "!");
                        staff.sendMessage(ChatColor.GRAY + "Link: " + ChatColor.AQUA + scanUrl);

                        if (!plugin.getScreenshareManager().isFrozen(target.getUniqueId())) {
                            plugin.getScreenshareManager().freezePlayer(target, staff);
                        }

                        for (String line : plugin.getConfig().getStringList("scan.link-message")) {
                            target.sendMessage(ChatColor.translateAlternateColorCodes('&',
                                line.replace("{link}", scanUrl)));
                        }
                    });
                } else {
                    Bukkit.getScheduler().runTask(plugin, () ->
                        staff.sendMessage(ChatColor.RED + "Failed to generate scan link. API error: " + response.statusCode()));
                }
            } catch (Exception e) {
                Bukkit.getScheduler().runTask(plugin, () ->
                    staff.sendMessage(ChatColor.RED + "Failed to connect to API: " + e.getMessage()));
            }
        });

        return true;
    }
}
