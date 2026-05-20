const { assessmentQuestions, assessmentReport } = require("../../../utils/mock-data");
const { ROUTES, redirect } = require("../../../utils/route");
const { getValue, setAssessmentSession, setValue } = require("../../../utils/storage");

function normalizeAnswer(value) {
  return String(value || "").trim().toLowerCase();
}

Page({
  data: {
    questions: assessmentQuestions,
    visibleQuestions: assessmentQuestions.slice(0, 3),
    currentIndex: 0,
    currentQuestion: assessmentQuestions[0],
    currentAnswer: "",
    answers: {
      q1: "B",
      q3: "false"
    },
    indexText: "01",
    displayTotal: 20,
    progressPercent: 5,
    timerDisplay: "35:20",
    isFirst: true,
    isLast: false,
    source: "entry"
  },

  onLoad(query) {
    const session = getValue("assessmentSession", {});
    this.remainingSeconds = Number(session.remainingSeconds || 35 * 60 + 20);
    this.setData({
      source: query.source || session.source || "entry",
      answers: Object.assign({}, this.data.answers, session.answers || {})
    });
    this.syncCurrent();
    this.startTimer();
  },

  onUnload() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  },

  startTimer() {
    this.updateTimerText();
    this.timer = setInterval(() => {
      this.remainingSeconds = Math.max(0, this.remainingSeconds - 1);
      this.updateTimerText();
    }, 1000);
  },

  updateTimerText() {
    const minutes = Math.floor(this.remainingSeconds / 60);
    const seconds = String(this.remainingSeconds % 60).padStart(2, "0");
    this.setData({
      timerDisplay: `${minutes}:${seconds}`
    });
  },

  syncCurrent() {
    const currentQuestion = this.data.questions[this.data.currentIndex];
    const currentAnswer = this.data.answers[currentQuestion.id] || "";
    const visibleQuestions = this.data.questions
      .slice(this.data.currentIndex, this.data.currentIndex + 3)
      .map((question) => Object.assign({}, question, {
        userAnswer: this.data.answers[question.id] || ""
      }));
    const progressPercent = Math.round(((this.data.currentIndex + 1) / this.data.displayTotal) * 100);
    this.setData({
      currentQuestion,
      visibleQuestions,
      currentAnswer,
      indexText: String(this.data.currentIndex + 1).padStart(2, "0"),
      progressPercent,
      isFirst: this.data.currentIndex === 0,
      isLast: this.data.currentIndex === this.data.questions.length - 1
    });
  },

  saveQuestionAnswer(questionId, value) {
    const nextAnswers = Object.assign({}, this.data.answers, {
      [questionId]: value
    });
    const visibleQuestions = this.data.visibleQuestions.map((question) => {
      return question.id === questionId ? Object.assign({}, question, { userAnswer: value }) : question;
    });
    this.setData({
      answers: nextAnswers,
      visibleQuestions,
      currentAnswer: questionId === this.data.currentQuestion.id ? value : this.data.currentAnswer
    });
    setAssessmentSession({
      source: this.data.source,
      questionSetId: "stitch-ai-demo",
      answers: nextAnswers,
      remainingSeconds: this.remainingSeconds
    });
  },

  selectOption(event) {
    this.saveQuestionAnswer(event.currentTarget.dataset.qid, event.currentTarget.dataset.value);
  },

  selectJudge(event) {
    this.saveQuestionAnswer(event.currentTarget.dataset.qid, event.currentTarget.dataset.value);
  },

  handleFillInput(event) {
    this.saveQuestionAnswer(event.currentTarget.dataset.qid, event.detail.value);
  },

  goPrev() {
    if (this.data.isFirst) {
      wx.showToast({
        title: "已经是第一题",
        icon: "none"
      });
      return;
    }
    this.setData({
      currentIndex: this.data.currentIndex - 1
    });
    this.syncCurrent();
  },

  goNext() {
    if (!normalizeAnswer(this.data.answers[this.data.currentQuestion.id])) {
      wx.showToast({
        title: "请先完成本题",
        icon: "none"
      });
      return;
    }

    if (this.data.isLast) {
      this.completeAssessment();
      return;
    }

    this.setData({
      currentIndex: this.data.currentIndex + 1
    });
    this.syncCurrent();
  },

  completeAssessment() {
    const answers = this.data.answers;
    const correctCount = this.data.questions.reduce((count, question) => {
      return count + (normalizeAnswer(answers[question.id]) === normalizeAnswer(question.answer) ? 1 : 0);
    }, 0);
    const score = Math.round((correctCount / this.data.questions.length) * 100);
    const report = Object.assign({}, assessmentReport, {
      score,
      correct: correctCount,
      wrong: this.data.questions.length - correctCount,
      total: this.data.questions.length,
      answers,
      source: this.data.source,
      completedAt: Date.now()
    });
    setValue("assessmentReport", report);
    setAssessmentSession({
      source: this.data.source,
      questionSetId: "stitch-ai-demo",
      answers,
      elapsedSeconds: 35 * 60 + 20 - this.remainingSeconds,
      reportId: "local-demo"
    });
    redirect(ROUTES.assessmentReport);
  }
});
