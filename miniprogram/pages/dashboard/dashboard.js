const { dashboard } = require("../../utils/mock-data");
const { ROUTES, navigate } = require("../../utils/route");
const { setAssessmentSession, setLearningContext } = require("../../utils/storage");

Page({
  data: {
    dashboard
  },

  handleToolTap(event) {
    const key = event.currentTarget.dataset.key;
    if (key === "short-test") {
      this.startWeaknessPractice("dashboard");
      return;
    }

    if (key === "study") {
      setLearningContext({
        subject: "英语",
        unit: "unit2",
        mode: "recognition"
      });
      navigate(ROUTES.studySetup);
      return;
    }

    if (key === "photo") {
      wx.showToast({
        title: "拍照解析即将开放",
        icon: "none"
      });
      return;
    }

    wx.showToast({
      title: "导师答疑即将开放",
      icon: "none"
    });
  },

  openRecommend(event) {
    const id = event.currentTarget.dataset.id;
    const mode = id === "english-cloze" ? "recognition" : "speaking";
    setLearningContext({
      gradeStage: "初中 (7-9)",
      subject: id === "english-cloze" ? "英语" : "数学",
      unit: id === "english-cloze" ? "unit2" : "unit1",
      mode
    });
    navigate(ROUTES.studySetup, { from: "recommend", id });
  },

  openKnowledgeAction(event) {
    const key = event.currentTarget.dataset.key;
    if (key === "practice") {
      this.startWeaknessPractice("knowledge-plan");
      return;
    }

    setLearningContext({
      gradeStage: "初中 (7-9)",
      subject: "英语",
      unit: "unit2",
      mode: key === "speaking" ? "speaking" : "recognition"
    });
    navigate(ROUTES.studySetup, { from: "knowledge-plan", mode: key });
  },

  openCourse(event) {
    const id = event.currentTarget.dataset.id;
    const course = this.data.dashboard.courses.find((item) => item.id === id);
    if (!course) {
      return;
    }

    if (course.mode === "practice") {
      this.startWeaknessPractice("course");
      return;
    }

    setLearningContext({
      gradeStage: "初中 (7-9)",
      subject: "英语",
      unit: course.unit,
      mode: course.mode
    });
    navigate(ROUTES.studySetup, { from: "course", id: course.id });
  },

  startWeaknessPractice(source) {
    setAssessmentSession({
      source,
      questionSetId: "weakness-demo",
      answers: {},
      elapsedSeconds: 0
    });
    navigate(ROUTES.exam, { source });
  },

  askAssistant() {
    wx.showToast({
      title: "AI 助教即将开放",
      icon: "none"
    });
  },

  openKnowledge() {
    navigate(ROUTES.knowledge);
  }
});
