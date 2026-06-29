package net.swiftsmp.anticheat.managers;

import net.swiftsmp.anticheat.SwiftACPlugin;
import org.bukkit.Bukkit;

import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.Scanner;
import java.util.concurrent.CompletableFuture;

public class APIManager {
    private final SwiftACPlugin plugin;
    private final String baseUrl;
    private final String apiKey;

    public APIManager(SwiftACPlugin plugin) {
        this.plugin = plugin;
        this.baseUrl = plugin.getConfig().getString("api.base-url", "http://localhost:3000");
        this.apiKey = plugin.getConfig().getString("api.api-key", "");
    }

    public CompletableFuture<String> createScanLink(String playerName, String playerUUID, String staffName) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                URI uri = new URI(baseUrl + "/api/scans/create");
                HttpURLConnection conn = (HttpURLConnection) uri.toURL().openConnection();
                conn.setRequestMethod("POST");
                conn.setRequestProperty("Content-Type", "application/json");
                conn.setRequestProperty("Authorization", "Bearer " + apiKey);
                conn.setDoOutput(true);

                String json = String.format(
                    "{\"playerName\":\"%s\",\"playerUUID\":\"%s\",\"staffName\":\"%s\"}",
                    playerName, playerUUID, staffName
                );

                try (OutputStream os = conn.getOutputStream()) {
                    os.write(json.getBytes(StandardCharsets.UTF_8));
                }

                int code = conn.getResponseCode();
                if (code == 200 || code == 201) {
                    try (Scanner s = new Scanner(conn.getInputStream(), StandardCharsets.UTF_8)) {
                        return s.useDelimiter("\\A").hasNext() ? s.next() : null;
                    }
                }
            } catch (Exception e) {
                plugin.getLogger().warning("API request failed: " + e.getMessage());
            }
            return null;
        });
    }
}
