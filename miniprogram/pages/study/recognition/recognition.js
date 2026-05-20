const { recognitionWords } = require("../../../utils/mock-data");
const { ROUTES, redirect } = require("../../../utils/route");
const { setStudySession } = require("../../../utils/storage");

Page({
  data: {
    words: recognitionWords,
    currentIndex: 0,
    current: recognitionWords[0],
    displayIndex: 5,
    displayTotal: 20,
    progressPercent: 25,
    showMeaning: true,
    knownCount: 0,
    fuzzyCount: 0
  },

  onLoad() {
    this.syncCurrent();
  },

  syncCurrent() {
    const current = this.data.words[this.data.currentIndex];
    const displayIndex = this.data.currentIndex + 5;
    this.setData({
      current,
      displayIndex,
      progressPercent: Math.round((displayIndex / this.data.displayTotal) * 100)
    });
  },

  playWord() {
    wx.showToast({
      title: `播放 ${this.data.current.word}`,
      icon: "none"
    });
  },

  toggleMeaning() {
    this.setData({
      showMeaning: !this.data.showMeaning
    });
  },

  openHelper(event) {
    wx.showToast({
      title: event.currentTarget.dataset.name,
      icon: "none"
    });
  },

  markFuzzy() {
    this.moveNext({
      fuzzyCount: this.data.fuzzyCount + 1
    });
  },

  markKnown() {
    this.moveNext({
      knownCount: this.data.knownCount + 1
    });
  },

  moveNext(patch) {
    const nextData = Object.assign({}, patch);
    if (this.data.currentIndex >= this.data.words.length - 1) {
      const knownCount = patch.knownCount === undefined ? this.data.knownCount : patch.knownCount;
      const fuzzyCount = patch.fuzzyCount === undefined ? this.data.fuzzyCount : patch.fuzzyCount;
      setStudySession({
        mode: "recognition",
        knownCount,
        fuzzyCount,
        total: this.data.words.length,
        score: Math.round((knownCount / this.data.words.length) * 100)
      });
      redirect(ROUTES.studyReport, { mode: "recognition" });
      return;
    }

    nextData.currentIndex = this.data.currentIndex + 1;
    nextData.showMeaning = true;
    this.setData(nextData);
    this.syncCurrent();
  }
});
