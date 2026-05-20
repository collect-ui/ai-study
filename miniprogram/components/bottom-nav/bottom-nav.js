Component({
  properties: {
    selected: {
      type: String,
      value: "home"
    },
    variant: {
      type: String,
      value: "primary"
    }
  },

  data: {
    items: []
  },

  lifetimes: {
    attached() {
      this.syncItems();
    }
  },

  observers: {
    "variant": function () {
      this.syncItems();
    }
  },

  methods: {
    syncItems() {
      const primaryItems = [
        { key: "home", label: "首页", icon: "⌂", path: "/pages/dashboard/dashboard" },
        { key: "profile", label: "我的", icon: "◎", path: "/pages/profile/english/english" }
      ];
      const fullItems = [
        { key: "home", label: "首页", icon: "⌂", path: "/pages/dashboard/dashboard" },
        { key: "study", label: "学习", icon: "□", path: "/pages/study/setup/setup" },
        { key: "knowledge", label: "知识", icon: "▥", path: "/pages/knowledge/overview/overview" },
        { key: "profile", label: "我的", icon: "◎", path: "/pages/profile/english/english" }
      ];

      this.setData({
        items: this.properties.variant === "full" ? fullItems : primaryItems
      });
    },

    navTo(event) {
      const { key, path } = event.currentTarget.dataset;
      this.triggerEvent("nav", { key, path });
      if (key === this.properties.selected) {
        return;
      }
      wx.redirectTo({
        url: path,
        fail: () => {
          wx.reLaunch({ url: path });
        }
      });
    }
  }
});
