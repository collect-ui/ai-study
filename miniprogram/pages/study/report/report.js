const { studyReport } = require("../../../utils/mock-data");
const { ROUTES, navigate, relaunch } = require("../../../utils/route");
const { getValue } = require("../../../utils/storage");

Page({
  data: {
    report: studyReport,
    mode: "speaking",
    playingAll: false,
    playingId: ""
  },

  onLoad(query) {
    const session = getValue("studySession", {});
    const report = Object.assign({}, studyReport);
    if (session.mode === "recognition") {
      report.score = Math.max(60, session.score || studyReport.score);
      report.comment = `已完成 ${session.total || 5} 个词卡，认识 ${session.knownCount || 0} 个`;
      report.metrics = [
        { label: "认识词卡", value: `${session.knownCount || 0}个` },
        { label: "模糊词卡", value: `${session.fuzzyCount || 0}个` },
        { label: "完成数量", value: `${session.total || 5}个` }
      ];
    } else if (session.completed && session.completed.length) {
      report.score = session.score || studyReport.score;
      report.details = session.completed;
    }

    this.setData({
      report,
      mode: query.mode || session.mode || "speaking"
    });
  },

  playAll() {
    this.setData({
      playingAll: !this.data.playingAll,
      playingId: ""
    });
    wx.showToast({
      title: this.data.playingAll ? "开始播放" : "已暂停播放",
      icon: "none"
    });
  },

  playItem(event) {
    const id = event.currentTarget.dataset.id;
    this.setData({
      playingId: id,
      playingAll: false
    });
    wx.showToast({
      title: "播放单句录音",
      icon: "none"
    });
  },

  backHome() {
    relaunch(ROUTES.dashboard);
  },

  nextGroup() {
    navigate(ROUTES.studySetup, { next: 1 });
  }
});
