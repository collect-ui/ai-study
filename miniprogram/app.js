App({
  globalData: {
    appName: "AI Study"
  },

  onLaunch() {
    const profile = wx.getStorageSync("profile");
    if (!profile) {
      wx.setStorageSync("profile", {
        streak: 0,
        learnedWords: 0,
        targetMinutes: 20
      });
    }

    const session = wx.getStorageSync("session");
    if (!session) {
      wx.setStorageSync("session", {
        isLoggedIn: false,
        isGuest: false
      });
    }

    const learningContext = wx.getStorageSync("learningContext");
    if (!learningContext) {
      wx.setStorageSync("learningContext", {
        gradeStage: "初中",
        grade: "初一",
        subject: "英语",
        unit: "unit1",
        mode: "recognition"
      });
    }
  }
});
