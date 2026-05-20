const { profile } = require("../../../utils/mock-data");
const { ROUTES, navigate } = require("../../../utils/route");
const { setLearningContext } = require("../../../utils/storage");

Page({
  data: {
    profile
  },

  openKnowledge() {
    navigate(ROUTES.knowledge, { from: "profile" });
  },

  openStudy() {
    setLearningContext({
      subject: "英语",
      unit: "unit1",
      mode: "recognition"
    });
    navigate(ROUTES.studySetup, { from: "profile" });
  },

  openMistakes() {
    navigate(ROUTES.mistakes, { from: "profile" });
  }
});
