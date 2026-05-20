const { ASSETS } = require("../../config/assets");

Component({
  properties: {
    title: {
      type: String,
      value: ""
    },
    subtitle: {
      type: String,
      value: ""
    },
    brand: {
      type: Boolean,
      value: false
    },
    showBack: {
      type: Boolean,
      value: false
    },
    showSearch: {
      type: Boolean,
      value: false
    },
    showBell: {
      type: Boolean,
      value: false
    },
    showMore: {
      type: Boolean,
      value: false
    },
    showShare: {
      type: Boolean,
      value: false
    },
    showSettings: {
      type: Boolean,
      value: false
    },
    showProfile: {
      type: Boolean,
      value: false
    },
    compact: {
      type: Boolean,
      value: false
    }
  },

  data: {
    assets: ASSETS
  },

  methods: {
    goBack() {
      this.triggerEvent("back");
      const pages = getCurrentPages();
      if (pages.length > 1) {
        wx.navigateBack();
        return;
      }
      wx.reLaunch({
        url: "/pages/dashboard/dashboard"
      });
    },

    emitAction(event) {
      const action = event.currentTarget.dataset.action;
      this.triggerEvent(action);
      if (action === "profile") {
        wx.navigateTo({
          url: "/pages/profile/english/english"
        });
        return;
      }
      const fallback = {
        search: "搜索功能待接入",
        bell: "暂无新通知",
        more: "更多操作待接入",
        share: "分享面板待接入",
        settings: "设置功能待接入",
        profile: "个人中心"
      };
      if (fallback[action]) {
        wx.showToast({
          title: fallback[action],
          icon: "none"
        });
      }
    }
  }
});
