const words = [
  {
    word: "focus",
    phonetic: "/ˈfoʊkəs/",
    meaning: "专注；焦点",
    example: "Set one clear goal and focus on it for twenty minutes.",
    tag: "学习习惯"
  },
  {
    word: "review",
    phonetic: "/rɪˈvjuː/",
    meaning: "复习；回顾",
    example: "A short review after class helps memory become stronger.",
    tag: "记忆"
  },
  {
    word: "practice",
    phonetic: "/ˈpræktɪs/",
    meaning: "练习；实践",
    example: "Daily practice is more useful than a long study once a week.",
    tag: "训练"
  }
];

const plan = [
  {
    title: "词汇热身",
    minutes: 5,
    detail: "快速浏览今日单词和例句"
  },
  {
    title: "理解训练",
    minutes: 10,
    detail: "跟读例句，理解真实语境"
  },
  {
    title: "复盘打卡",
    minutes: 5,
    detail: "记录今日完成情况"
  }
];

module.exports = {
  words,
  plan
};
