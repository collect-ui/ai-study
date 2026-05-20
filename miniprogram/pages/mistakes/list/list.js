const { mistakeCategories, mistakeItems } = require("../../../utils/mock-data");
const { ROUTES, navigate } = require("../../../utils/route");
const { setAssessmentSession } = require("../../../utils/storage");

Page({
  data: {
    categories: mistakeCategories,
    items: mistakeItems,
    filtered: [],
    activeCategory: "words",
    keyword: "",
    activeAnalysisId: "",
    stats: {
      total: 128,
      todo: 12,
      mastered: 86,
      percent: 67
    }
  },

  onLoad() {
    this.refreshList();
  },

  refreshList() {
    const keyword = this.data.keyword.trim().toLowerCase();
    const prototypeIds = ["m1", "m2", "m4"];
    const filtered = this.data.items.filter((item) => {
      const categoryMatched = !keyword && this.data.activeCategory === "words"
        ? prototypeIds.indexOf(item.id) >= 0
        : item.category === this.data.activeCategory;
      const keywordMatched = !keyword || item.title.toLowerCase().indexOf(keyword) >= 0 || item.tag.toLowerCase().indexOf(keyword) >= 0;
      return categoryMatched && keywordMatched;
    });
    this.setData({ filtered });
  },

  selectCategory(event) {
    this.setData({
      activeCategory: event.currentTarget.dataset.id,
      activeAnalysisId: ""
    });
    this.refreshList();
  },

  handleSearchInput(event) {
    this.setData({
      keyword: event.detail.value
    });
    this.refreshList();
  },

  toggleAnalysis(event) {
    const id = event.currentTarget.dataset.id;
    this.setData({
      activeAnalysisId: this.data.activeAnalysisId === id ? "" : id
    });
  },

  startReview() {
    setAssessmentSession({
      source: "mistake",
      questionSetId: "mistake-review-demo",
      answers: {},
      elapsedSeconds: 0
    });
    navigate(ROUTES.exam, { source: "mistake" });
  }
});
