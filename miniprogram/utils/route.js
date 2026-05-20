const ROUTES = {
  home: "/pages/home/home",
  dashboard: "/pages/dashboard/dashboard",
  exam: "/pages/assessment/exam/exam",
  assessmentReport: "/pages/assessment/report/report",
  studySetup: "/pages/study/setup/setup",
  recognition: "/pages/study/recognition/recognition",
  speaking: "/pages/study/speaking/speaking",
  studyReport: "/pages/study/report/report",
  knowledge: "/pages/knowledge/overview/overview",
  mistakes: "/pages/mistakes/list/list",
  profile: "/pages/profile/english/english"
};

function buildUrl(path, query) {
  if (!query) {
    return path;
  }
  const pairs = Object.keys(query)
    .filter((key) => query[key] !== undefined && query[key] !== null && query[key] !== "")
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(query[key])}`);
  return pairs.length ? `${path}?${pairs.join("&")}` : path;
}

function navigate(path, query) {
  wx.navigateTo({
    url: buildUrl(path, query)
  });
}

function redirect(path, query) {
  wx.redirectTo({
    url: buildUrl(path, query)
  });
}

function relaunch(path, query) {
  wx.reLaunch({
    url: buildUrl(path, query)
  });
}

function goBackOrHome() {
  const pages = getCurrentPages();
  if (pages.length > 1) {
    wx.navigateBack();
    return;
  }
  relaunch(ROUTES.dashboard);
}

module.exports = {
  ROUTES,
  buildUrl,
  navigate,
  redirect,
  relaunch,
  goBackOrHome
};
