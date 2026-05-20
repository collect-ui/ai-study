const { ROUTES, navigate, relaunch } = require("../../utils/route");
const { setLearningContext, setSession } = require("../../utils/storage");
const { ASSETS } = require("../../config/assets");

Page({
  data: {
    phone: "",
    code: "",
    canSendCode: true,
    codeButtonText: "获取验证码",
    stages: ["小学", "初中"],
    selectedStage: "小学",
    allGrades: {
      "小学": ["一年级", "二年级", "三年级", "四年级", "五年级", "六年级"],
      "初中": ["初一", "初二", "初三"]
    },
    grades: ["一年级", "二年级", "三年级", "四年级", "五年级", "六年级"],
    selectedGrade: "三年级",
    assets: ASSETS,
    subjects: [
      { id: "chinese", name: "语文", icon: ASSETS.book },
      { id: "math", name: "数学", icon: ASSETS.sigma },
      { id: "english", name: "英语", icon: ASSETS.translate }
    ],
    selectedSubject: "chinese"
  },

  handlePhoneInput(event) {
    this.setData({
      phone: event.detail.value
    });
  },

  handleCodeInput(event) {
    this.setData({
      code: event.detail.value
    });
  },

  sendCode() {
    if (!this.data.canSendCode) {
      return;
    }

    if (!/^1\d{10}$/.test(this.data.phone)) {
      wx.showToast({
        title: "请输入正确手机号",
        icon: "none"
      });
      return;
    }

    this.setData({
      canSendCode: false,
      codeButtonText: "已发送 60s"
    });

    wx.showToast({
      title: "验证码已发送",
      icon: "none"
    });
  },

  login() {
    if (!/^1\d{10}$/.test(this.data.phone)) {
      wx.showToast({
        title: "请输入正确手机号",
        icon: "none"
      });
      return;
    }

    if (!/^\d{4,6}$/.test(this.data.code)) {
      wx.showToast({
        title: "请输入验证码",
        icon: "none"
      });
      return;
    }

    this.persistLearningContext();
    setSession({
      isLoggedIn: true,
      isGuest: false,
      profileId: "student-demo"
    });

    wx.showToast({
      title: "登录成功",
      icon: "success",
      duration: 700
    });

    setTimeout(() => {
      relaunch(ROUTES.dashboard);
    }, 700);
  },

  continueAsGuest() {
    this.persistLearningContext();
    setSession({
      isLoggedIn: false,
      isGuest: true,
      profileId: "guest-demo"
    });

    wx.showToast({
      title: "已进入游客模式",
      icon: "none"
    });

    setTimeout(() => {
      relaunch(ROUTES.dashboard);
    }, 500);
  },

  showNotifications() {
    wx.showToast({
      title: "暂无通知",
      icon: "none"
    });
  },

  selectStage(event) {
    const selectedStage = event.currentTarget.dataset.stage;
    const grades = this.data.allGrades[selectedStage] || [];
    this.setData({
      selectedStage,
      grades,
      selectedGrade: grades[0] || ""
    });
  },

  selectGrade(event) {
    this.setData({
      selectedGrade: event.currentTarget.dataset.grade
    });
  },

  selectSubject(event) {
    this.setData({
      selectedSubject: event.currentTarget.dataset.subject
    });
  },

  startAssessment() {
    const subject = this.data.subjects.find((item) => item.id === this.data.selectedSubject);
    this.persistLearningContext();
    navigate(ROUTES.exam, {
      source: "entry",
      grade: this.data.selectedGrade,
      subject: subject ? subject.name : ""
    });
  },

  persistLearningContext() {
    const subject = this.data.subjects.find((item) => item.id === this.data.selectedSubject);
    setLearningContext({
      gradeStage: this.data.selectedStage,
      grade: this.data.selectedGrade,
      subject: subject ? subject.name : "语文"
    });
  }
});
