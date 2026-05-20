const { words } = require("../../utils/study-data");

Page({
  data: {
    words,
    currentIndex: 0,
    current: words[0]
  },

  nextWord() {
    const nextIndex = (this.data.currentIndex + 1) % this.data.words.length;
    this.setData({
      currentIndex: nextIndex,
      current: this.data.words[nextIndex]
    });
  },

  markLearned() {
    const profile = wx.getStorageSync("profile") || {};
    const learnedWords = Number(profile.learnedWords || 0) + 1;
    wx.setStorageSync("profile", {
      streak: Number(profile.streak || 0),
      targetMinutes: Number(profile.targetMinutes || 20),
      learnedWords
    });
    wx.showToast({
      title: "已记录",
      icon: "success"
    });
  },

  playPronunciation() {
    wx.showToast({
      title: "后续接入发音",
      icon: "none"
    });
  }
});
