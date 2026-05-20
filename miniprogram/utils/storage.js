function getValue(key, fallback) {
  const value = wx.getStorageSync(key);
  return value === "" || value === undefined || value === null ? fallback : value;
}

function setValue(key, value) {
  wx.setStorageSync(key, value);
  return value;
}

function mergeValue(key, patch) {
  const current = getValue(key, {});
  const next = Object.assign({}, current, patch);
  wx.setStorageSync(key, next);
  return next;
}

function setSession(patch) {
  return mergeValue("session", patch);
}

function setLearningContext(patch) {
  return mergeValue("learningContext", patch);
}

function setAssessmentSession(value) {
  return setValue("assessmentSession", value);
}

function setStudySession(value) {
  return setValue("studySession", value);
}

module.exports = {
  getValue,
  setValue,
  mergeValue,
  setSession,
  setLearningContext,
  setAssessmentSession,
  setStudySession
};
