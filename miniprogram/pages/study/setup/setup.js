const { studySetup } = require("../../../utils/mock-data");
const { ROUTES, navigate } = require("../../../utils/route");
const { getValue, setLearningContext } = require("../../../utils/storage");

Page({
  data: {
    setup: studySetup,
    selectedGrade: "junior",
    selectedUnit: "unit1",
    selectedMode: "recognition"
  },

  onLoad() {
    const context = getValue("learningContext", {});
    this.setData({
      selectedUnit: context.unit || "unit1",
      selectedMode: context.mode || "recognition"
    });
  },

  selectGrade(event) {
    this.setData({
      selectedGrade: event.currentTarget.dataset.id
    });
  },

  selectUnit(event) {
    this.setData({
      selectedUnit: event.currentTarget.dataset.id
    });
  },

  selectMode(event) {
    this.setData({
      selectedMode: event.currentTarget.dataset.id
    });
  },

  startStudy() {
    const grade = this.data.setup.grades.find((item) => item.id === this.data.selectedGrade);
    const unit = this.data.setup.units.find((item) => item.id === this.data.selectedUnit);
    setLearningContext({
      gradeStage: grade ? grade.name : "初中 (7-9)",
      subject: this.data.setup.subject.name,
      unit: this.data.selectedUnit,
      unitTitle: unit ? unit.title : "",
      mode: this.data.selectedMode
    });

    navigate(this.data.selectedMode === "speaking" ? ROUTES.speaking : ROUTES.recognition);
  }
});
