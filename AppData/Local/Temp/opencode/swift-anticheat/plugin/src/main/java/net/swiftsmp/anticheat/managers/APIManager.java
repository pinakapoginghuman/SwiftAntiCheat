package net.swiftsmp.anticheat.managers;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import net.swiftsmp.anticheat.SwiftACPlugin;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

public class APIManager {
    private final SwiftACPlugin plugin;
    private final String baseUrl;
    private final String apiKey;
    private final HttpClient client;
    private final Gson gson;

    public APIManager(SwiftACPlugin plugin) {
        this.plugin = plugin;
        this.baseUrl = plugin.getConfig().getString("api.base-url", "http://localhost:3000");
        this.apiKey = plugin.getConfig().getString("api.api-key", "");
        this.client = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(10))
            .build();
        this.gson = new Gson();
    }

    public ScanResponse createScan(String playerName, String playerUUID, String staffName) {
        try {
            JsonObject body = new JsonObject();
            body.addProperty("playerName", playerName);
            body.addProperty("playerUUID", playerUUID);
            body.addProperty("staffName", staffName);

            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/api/scans/create"))
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer " + apiKey)
                .POST(HttpRequest.BodyPublishers.ofString(gson.toJson(body)))
                .timeout(Duration.ofSeconds(10))
                .build();

            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

            if (response.statusCode() == 200 || response.statusCode() == 201) {
                JsonObject json = gson.fromJson(response.body(), JsonObject.class);
                return new ScanResponse(
                    json.get("id").getAsString(),
                    json.get("status").getAsString()
                );
            } else {
                plugin.getLogger().warning("API create scan failed: " + response.statusCode() + " " + response.body());
            }
        } catch (Exception e) {
            plugin.getLogger().warning("API request failed: " + e.getMessage());
        }
        return null;
    }

    public static class ScanResponse {
        private final String id;
        private final String status;

        public ScanResponse(String id, String status) {
            this.id = id;
            this.status = status;
        }

        public String getId() { return id; }
        public String getStatus() { return status; }
    }
}
