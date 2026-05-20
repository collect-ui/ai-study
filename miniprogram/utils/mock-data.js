const dashboard = {
  studentName: "小明",
  grade: "初三",
  subject: "英语",
  achievement: {
    percent: 75,
    remain: 3,
    stages: ["基础巩固", "查漏补缺", "冲刺提升"]
  },
  progress: {
    term: "九年级上 · Unit 2",
    weekly: "本周已学 4/6 课",
    next: "今天建议学习 25 分钟",
    skills: [
      { name: "重点句型", percent: 62, status: "待巩固" },
      { name: "阅读理解", percent: 74, status: "继续练" },
      { name: "听说跟读", percent: 58, status: "需加强" }
    ]
  },
  knowledgePlan: [
    {
      id: "sentence",
      title: "but / because 连接句",
      desc: "先完成 8 道句型辨析，解决最近错题。",
      action: "去练习",
      actionKey: "practice",
      icon: "句",
      tone: "blue"
    },
    {
      id: "vocab",
      title: "Unit 2 高频词汇",
      desc: "friendship、practice、advice 今天需要认读。",
      action: "认读",
      actionKey: "recognition",
      icon: "词",
      tone: "green"
    },
    {
      id: "oral",
      title: "每日跟读 4 句",
      desc: "重点纠正 perfect / catches 的发音。",
      action: "跟读",
      actionKey: "speaking",
      icon: "读",
      tone: "gold"
    }
  ],
  courses: [
    {
      id: "unit1",
      title: "Unit 1: Making New Friends",
      meta: "词汇认读 · 已完成 80%",
      progress: 80,
      unit: "unit1",
      mode: "recognition",
      action: "继续"
    },
    {
      id: "unit2",
      title: "Unit 2: My School Life",
      meta: "重点句型 · 当前学习",
      progress: 45,
      unit: "unit2",
      mode: "speaking",
      action: "开始"
    },
    {
      id: "grammar",
      title: "时态专项：过去进行时",
      meta: "语法练习 · 12 分钟",
      progress: 20,
      unit: "unit2",
      mode: "practice",
      action: "练习"
    }
  ],
  tools: [
    {
      id: "short-test",
      title: "短板测试",
      desc: "精准锁定知识盲区",
      icon: "测",
      tag: "AI 诊断",
      tone: "red",
      wide: true
    },
    {
      id: "photo",
      title: "拍照录入",
      desc: "秒级识别错题与解析",
      icon: "拍",
      tone: "green"
    },
    {
      id: "study",
      title: "自主学习",
      desc: "名师精讲微课合集",
      icon: "学",
      tone: "blue"
    },
    {
      id: "qa",
      title: "在线答疑",
      desc: "真人在校导师 1对1 实时解惑",
      icon: "问",
      tone: "gold",
      tag: "立即连接",
      wide: true
    }
  ],
  recommended: [
    {
      id: "math-geometry",
      subject: "数学 · 解析几何",
      title: "三步搞定椭圆最值问题",
      desc: "1,240 位同学正在学习",
      tone: "blue",
      icon: "数"
    },
    {
      id: "english-cloze",
      subject: "英语 · 完型填空",
      title: "高频词汇智能闪卡训练",
      desc: "基于你昨天的错题生成",
      tone: "green",
      icon: "词"
    }
  ]
};

const assessmentQuestions = [
  {
    id: "q1",
    type: "single",
    typeLabel: "单选题",
    helper: "请根据句意选择最恰当的单词填空。",
    question: "The Great Wall is one of the _____ wonders in the world.",
    answer: "B",
    options: [
      { value: "A", text: "great" },
      { value: "B", text: "greatest" },
      { value: "C", text: "greater" },
      { value: "D", text: "greatly" }
    ],
    analysis: "one of the 后接最高级和名词复数，great 的最高级是 greatest。"
  },
  {
    id: "q2",
    type: "fill",
    typeLabel: "填空题",
    helper: "写出括号内动词的正确形式。",
    question: "My mother (cook) dinner when I got home yesterday.",
    placeholder: "输入答案...",
    answer: "was cooking",
    analysis: "when 引导过去动作发生时，主句表示正在进行，用过去进行时。"
  },
  {
    id: "q3",
    type: "judge",
    typeLabel: "判断题",
    helper: "判断下列陈述是否正确。",
    passage: "Li Hua is a 14-year-old student. He likes sports very much. He plays basketball every afternoon. His favorite subject is English because he thinks it is very useful.",
    question: "Li Hua 最喜欢的科目是体育，因为他喜欢运动。",
    answer: "false",
    analysis: "原文说明他最喜欢的科目是 English，而不是体育。"
  },
  {
    id: "q4",
    type: "single",
    typeLabel: "单选题",
    helper: "选择能完成对话的正确表达。",
    question: "Could you tell me _____ the library is?",
    answer: "C",
    options: [
      { value: "A", text: "what" },
      { value: "B", text: "when" },
      { value: "C", text: "where" },
      { value: "D", text: "who" }
    ],
    analysis: "询问地点用 where。"
  },
  {
    id: "q5",
    type: "fill",
    typeLabel: "填空题",
    helper: "根据语境补全单词。",
    question: "A good friend always gives you helpful a_____.",
    placeholder: "输入单词...",
    answer: "advice",
    analysis: "helpful advice 表示有帮助的建议，advice 是不可数名词。"
  },
  {
    id: "q6",
    type: "judge",
    typeLabel: "判断题",
    helper: "判断语法说明是否正确。",
    passage: "Practice makes perfect. The sentence uses the simple present tense to express a general truth.",
    question: "这句话使用一般现在时表达普遍真理。",
    answer: "true",
    analysis: "一般现在时常用于表达客观事实或普遍真理。"
  }
];

const assessmentReport = {
  score: 85,
  level: "优秀",
  beatRate: 92,
  correct: 17,
  wrong: 3,
  mastery: [
    { name: "词汇语法", percent: 90 },
    { name: "阅读理解", percent: 75 },
    { name: "逻辑推理", percent: 82 }
  ],
  suggestion: "词汇语法基础扎实，阅读速度和长难句拆分仍有提升空间。建议优先复盘错题并完成一组针对性练习。"
};

const studySetup = {
  grades: [
    { id: "junior", name: "初中 (7-9)" },
    { id: "primary", name: "小学 (1-6)" },
    { id: "senior", name: "高中 (10-12)" }
  ],
  subject: {
    id: "english",
    name: "英语 (English)"
  },
  units: [
    { id: "unit1", title: "Unit 1: Making New Friends", icon: "1" },
    { id: "unit2", title: "Unit 2: My School Life", icon: "2" },
    { id: "unit3", title: "Unit 3: Food and Health", icon: "3" }
  ],
  modes: [
    { id: "recognition", title: "认读模式", desc: "看词识义，强化记忆", icon: "认" },
    { id: "speaking", title: "跟读模式", desc: "标准发音，纠正口语", icon: "读" }
  ]
};

const recognitionWords = [
  {
    word: "Friendship",
    phonetic: "/ˈfrend.ʃɪp/",
    meaning: "友谊",
    detail: "n. 友谊，友情，友好关系",
    last: "3天前",
    mastery: 25
  },
  {
    word: "Treasure",
    phonetic: "/ˈtreʒ.ər/",
    meaning: "珍宝",
    detail: "n. 宝物；v. 珍惜",
    last: "5天前",
    mastery: 42
  },
  {
    word: "Practice",
    phonetic: "/ˈpræk.tɪs/",
    meaning: "练习",
    detail: "n. 练习；v. 实践",
    last: "昨天",
    mastery: 63
  },
  {
    word: "Useful",
    phonetic: "/ˈjuːs.fəl/",
    meaning: "有用的",
    detail: "adj. 有帮助的，实用的",
    last: "今天",
    mastery: 76
  },
  {
    word: "Indeed",
    phonetic: "/ɪnˈdiːd/",
    meaning: "确实",
    detail: "adv. 的确，真正地",
    last: "7天前",
    mastery: 31
  }
];

const speakingSentences = [
  {
    text: "A friend in need is a friend indeed.",
    cn: "患难见真情。",
    focus: "need",
    score: 95,
    feedback: "非常流利",
    advice: "整体语调很好，单词 need 的长元音 /i:/ 可以再拉长。"
  },
  {
    text: "Practice makes perfect.",
    cn: "熟能生巧。",
    focus: "perfect",
    score: 89,
    feedback: "节奏稳定",
    advice: "注意 perfect 的重音落在第一音节。"
  },
  {
    text: "Better late than never.",
    cn: "迟做总比不做好。",
    focus: "better",
    score: 92,
    feedback: "发音清晰",
    advice: "连读处理自然，可以继续保持。"
  },
  {
    text: "The early bird catches the worm.",
    cn: "早起的鸟儿有虫吃。",
    focus: "catches",
    score: 86,
    feedback: "需要复听",
    advice: "catches 的词尾 /ɪz/ 不要吞音。"
  }
];

const studyReport = {
  score: 92,
  label: "优 Excellent",
  comment: "太棒了！你的发音非常标准",
  metrics: [
    { label: "正确率", value: "88%" },
    { label: "学习时长", value: "12:05" },
    { label: "练习句数", value: "10句" }
  ],
  details: [
    { sentence: "A friend in need is a friend indeed.", score: 95, feedback: "发音饱满，语调自然" },
    { sentence: "Practice makes perfect.", score: 89, feedback: "注意 perfect 的重音" },
    { sentence: "Better late than never.", score: 92, feedback: "连读处理得很棒" },
    { sentence: "The early bird catches the worm.", score: 86, feedback: "词尾爆破音较弱" }
  ]
};

const knowledgeStats = {
  score: 85,
  beatRate: 92,
  abilities: [
    { id: "vocabulary", name: "词汇", percent: 78, desc: "掌握 4,200 单词", tone: "blue", color: "#005da7" },
    { id: "grammar", name: "语法", percent: 92, desc: "核心语法点已精通", tone: "green", color: "#006b58" },
    { id: "oral", name: "口语表达", percent: 65, desc: "发音准确度高，需加强流利度练习与语调起伏感。", tone: "gold", color: "#a06900" }
  ],
  skills: [
    { name: "听力理解", percent: 82 },
    { name: "阅读能力", percent: 88 },
    { name: "写作表达", percent: 74 }
  ],
  insights: [
    "长难句分析能力有显著提升，但定语从句的灵活运用仍有薄弱点。",
    "建议本周加强口语模考，重点练习弱读与连读技巧。",
    "词汇记忆曲线显示，你需要复习 Unit 4-6 的科技类单词。"
  ]
};

const mistakeCategories = [
  { id: "words", name: "生字词" },
  { id: "sentences", name: "重点句子" },
  { id: "daily", name: "每日一题" },
  { id: "other", name: "其他" }
];

const mistakeItems = [
  {
    id: "m1",
    category: "words",
    tag: "高频错误",
    title: "请找出下列句子中的错别字：他不仅学习刻苦，而且待人诚恳。",
    date: "2023-10-24",
    mastery: 2,
    starText: "★★☆",
    analysis: "该题考查近义词和形近字辨析，复习时先读完整句再定位关键词。"
  },
  {
    id: "m2",
    category: "sentences",
    tag: "重点复习",
    title: "翻译句子：Knowledge is a treasure, but practice is the key to it.",
    date: "2023-10-22",
    mastery: 1,
    starText: "★☆☆",
    analysis: "句子使用转折结构 but，后半句 key to it 是固定搭配。"
  },
  {
    id: "m3",
    category: "daily",
    tag: "每日一题",
    title: "The Great Wall is one of the greatest wonders in the world.",
    date: "2023-10-21",
    mastery: 3,
    starText: "★★★",
    tutors: 2,
    analysis: "one of the + 最高级 + 名词复数，是中考常见结构。"
  },
  {
    id: "m4",
    category: "other",
    tag: "数学计算",
    title: "解方程：3x + 5 = 2x - 7",
    date: "2023-10-20",
    mastery: 2,
    starText: "★★☆",
    tutors: 2,
    analysis: "移项后 x = -12，注意符号变化。"
  }
];

const profile = {
  name: "陈小智",
  level: "Advanced Learner",
  grade: "高二年级 · 理科实验班",
  badge: "English Star",
  tags: ["单词王", "英语课代表"],
  mockScore: 138,
  mockTotal: 150,
  improvement: "+12%",
  tasks: {
    done: 8,
    total: 12,
    percent: 67
  },
  knowledge: [
    { name: "词汇 Vocabulary", percent: 85, tone: "green" },
    { name: "语法 Grammar", percent: 62, tone: "blue" },
    { name: "口语 Oral", percent: 45, tone: "gold" }
  ],
  mistakes: 88
};

module.exports = {
  dashboard,
  assessmentQuestions,
  assessmentReport,
  studySetup,
  recognitionWords,
  speakingSentences,
  studyReport,
  knowledgeStats,
  mistakeCategories,
  mistakeItems,
  profile
};
