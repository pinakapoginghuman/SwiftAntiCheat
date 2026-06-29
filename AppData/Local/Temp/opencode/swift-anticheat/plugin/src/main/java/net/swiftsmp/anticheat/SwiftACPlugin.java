package net.swiftsmp.anticheat;

import net.swiftsmp.anticheat.commands.ScreenshareCommand;
import net.swiftsmp.anticheat.commands.UnscreenshareCommand;
import net.swiftsmp.anticheat.commands.SwiftAntiCheatCommand;
import net.swiftsmp.anticheat.commands.DiscordCommand;
import net.swiftsmp.anticheat.listeners.PlayerListener;
import net.swiftsmp.anticheat.listeners.ChatListener;
import net.swiftsmp.anticheat.managers.ScreenshareManager;
import net.swiftsmp.anticheat.managers.APIManager;
import net.swiftsmp.anticheat.managers.BanManager;
import org.bukkit.plugin.java.JavaPlugin;

public class SwiftACPlugin extends JavaPlugin {
    private static SwiftACPlugin instance;
    private ScreenshareManager screenshareManager;
    private APIManager apiManager;
    private BanManager banManager;

    @Override
    public void onEnable() {
        instance = this;
        saveDefaultConfig();

        this.screenshareManager = new ScreenshareManager(this);
        this.apiManager = new APIManager(this);
        this.banManager = new BanManager(this);

        getCommand("screenshare").setExecutor(new ScreenshareCommand(this));
        getCommand("unscreenshare").setExecutor(new UnscreenshareCommand(this));
        getCommand("swiftanticheat").setExecutor(new SwiftAntiCheatCommand(this));
        getCommand("discord").setExecutor(new DiscordCommand(this));
        getServer().getPluginManager().registerEvents(new PlayerListener(this), this);
        getServer().getPluginManager().registerEvents(new ChatListener(this), this);

        getLogger().info("SwiftAntiCheat v" + getDescription().getVersion() + " enabled!");
    }

    @Override
    public void onDisable() {
        screenshareManager.unfreezeAll();
        getLogger().info("SwiftAntiCheat disabled.");
    }

    public static SwiftACPlugin getInstance() { return instance; }
    public ScreenshareManager getScreenshareManager() { return screenshareManager; }
    public APIManager getAPIManager() { return apiManager; }
    public BanManager getBanManager() { return banManager; }
}
