const { knowledgeStats } = require("../../../utils/mock-data");
const { ROUTES, navigate } = require("../../../utils/route");
const { setAssessmentSession } = require("../../../utils/storage");

Page({
  data: {
    stats: knowledgeStats
  },

  startPractice() {
    setAssessmentSession({
      source: "knowledge",
      questionSetId: "knowledge-target-demo",
      answers: {},
      elapsedSeconds: 0
    });
    navigate(ROUTES.exam, { source: "knowledge" });
  },

  openMistakes() {
    navigate(ROUTES.mistakes, { from: "knowledge" });
  }
});
