Page({
  data: {
    profile: {
      streak: 0,
      learnedWords: 0,
      targetMinutes: 20
    },
    checkedToday: false
  },

  onShow() {
    const today = this.getToday();
    const lastCheckin = wx.getStorageSync("lastCheckin");
    this.setData({
      profile: wx.getStorageSync("profile") || this.data.profile,
      checkedToday: lastCheckin === today
    });
  },

  checkin() {
    const today = this.getToday();
    if (this.data.checkedToday) {
      wx.showToast({
        title: "今天已打卡",
        icon: "none"
      });
      return;
    }

    const profile = wx.getStorageSync("profile") || this.data.profile;
    const nextProfile = {
      ...profile,
      streak: Number(profile.streak || 0) + 1
    };
    wx.setStorageSync("profile", nextProfile);
    wx.setStorageSync("lastCheckin", today);
    this.setData({
      profile: nextProfile,
      checkedToday: true
    });
    wx.showToast({
      title: "打卡成功",
      icon: "success"
    });
  },

  resetProgress() {
    const nextProfile = {
      streak: 0,
      learnedWords: 0,
      targetMinutes: 20
    };
    wx.setStorageSync("profile", nextProfile);
    wx.removeStorageSync("lastCheckin");
    this.setData({
      profile: nextProfile,
      checkedToday: false
    });
  },

  getToday() {
    const date = new Date();
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }
});
