const { speakingSentences } = require("../../../utils/mock-data");
const { ROUTES, redirect, navigate, relaunch } = require("../../../utils/route");
const { setStudySession } = require("../../../utils/storage");

Page({
  data: {
    sentences: speakingSentences,
    currentIndex: 0,
    current: speakingSentences[0],
    displayIndex: 3,
    displayTotal: 10,
    progressPercent: 30,
    recording: false,
    hasResult: true,
    score: speakingSentences[0].score,
    feedback: speakingSentences[0].feedback,
    advice: speakingSentences[0].advice,
    completed: []
  },

  onLoad() {
    this.syncCurrent();
  },

  syncCurrent() {
    const current = this.data.sentences[this.data.currentIndex];
    const displayIndex = this.data.currentIndex + 3;
    this.setData({
      current,
      displayIndex,
      progressPercent: Math.round((displayIndex / this.data.displayTotal) * 100)
    });
  },

  playSample() {
    wx.showToast({
      title: "播放标准范读",
      icon: "none"
    });
  },

  startRecord() {
    this.setData({
      recording: true,
      hasResult: false
    });
  },

  finishRecord() {
    if (!this.data.recording) {
      return;
    }
    this.applyResult();
  },

  tapMic() {
    if (this.data.hasResult) {
      wx.showToast({
        title: "长按可重新录音",
        icon: "none"
      });
      return;
    }
    this.applyResult();
  },

  applyResult() {
    this.setData({
      recording: false,
      hasResult: true,
      score: this.data.current.score,
      feedback: this.data.current.feedback,
      advice: this.data.current.advice
    });
    wx.showToast({
      title: "已生成评分",
      icon: "success"
    });
  },

  redo() {
    this.setData({
      recording: false,
      hasResult: false,
      score: 0,
      feedback: "等待录音",
      advice: "长按话筒完成本句跟读后，将生成 AI 建议。"
    });
  },

  nextSentence() {
    if (!this.data.hasResult) {
      wx.showToast({
        title: "请先完成录音",
        icon: "none"
      });
      return;
    }

    const completed = this.data.completed.concat({
      sentence: this.data.current.text,
      score: this.data.score,
      feedback: this.data.feedback
    });

    if (this.data.currentIndex >= this.data.sentences.length - 1) {
      setStudySession({
        mode: "speaking",
        score: Math.round(completed.reduce((sum, item) => sum + item.score, 0) / completed.length),
        completed
      });
      redirect(ROUTES.studyReport, { mode: "speaking" });
      return;
    }

    this.setData({
      currentIndex: this.data.currentIndex + 1,
      completed,
      hasResult: true,
      score: this.data.sentences[this.data.currentIndex + 1].score,
      feedback: this.data.sentences[this.data.currentIndex + 1].feedback,
      advice: this.data.sentences[this.data.currentIndex + 1].advice
    });
    this.syncCurrent();
  },

  goHome() {
    relaunch(ROUTES.dashboard);
  },

  goStudy() {
    navigate(ROUTES.studySetup);
  },

  goMistakes() {
    navigate(ROUTES.mistakes, { from: "speaking" });
  },

  goProfile() {
    navigate(ROUTES.profile);
  }
});
