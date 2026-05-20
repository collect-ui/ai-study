const { assessmentReport } = require("../../../utils/mock-data");
const { ASSETS } = require("../../../config/assets");
const { ROUTES, navigate } = require("../../../utils/route");
const { getValue, setAssessmentSession } = require("../../../utils/storage");

Page({
  data: {
    report: assessmentReport,
    assets: ASSETS,
    contactName: "",
    contactPhone: "",
    submitted: false
  },

  onLoad() {
    this.setData({
      report: getValue("assessmentReport", assessmentReport)
    });
  },

  handleNameInput(event) {
    this.setData({
      contactName: event.detail.value
    });
  },

  handlePhoneInput(event) {
    this.setData({
      contactPhone: event.detail.value
    });
  },

  submitConsult() {
    if (!this.data.contactName.trim()) {
      wx.showToast({
        title: "请输入姓名",
        icon: "none"
      });
      return;
    }

    if (!/^1\d{10}$/.test(this.data.contactPhone)) {
      wx.showToast({
        title: "请输入正确手机号",
        icon: "none"
      });
      return;
    }

    this.setData({
      submitted: true
    });
    wx.showToast({
      title: "已提交咨询",
      icon: "success"
    });
  },

  openMistakes() {
    navigate(ROUTES.mistakes, { from: "assessment" });
  },

  openKnowledge() {
    navigate(ROUTES.knowledge, { from: "assessment" });
  },

  makePractice() {
    setAssessmentSession({
      source: "report",
      questionSetId: "improvement-demo",
      answers: {},
      elapsedSeconds: 0
    });
    navigate(ROUTES.exam, { source: "report" });
  }
});
