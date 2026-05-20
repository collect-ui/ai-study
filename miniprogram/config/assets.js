const REMOTE_ASSET_BASE_URL = "https://collect-ui.top/ai-study/assets";
const REMOTE_ASSET_VERSION = "20260517-ca1";

function assetUrl(path) {
  return `${REMOTE_ASSET_BASE_URL}/${String(path).replace(/^\/+/, "")}?v=${REMOTE_ASSET_VERSION}`;
}

const ASSETS = {
  arrow: assetUrl("prototype/arrow.png"),
  bell: assetUrl("prototype/bell.png"),
  book: assetUrl("prototype/book.png"),
  logo: assetUrl("prototype/logo.png"),
  sigma: assetUrl("prototype/sigma.png"),
  sparkle: assetUrl("prototype/sparkle.png"),
  translate: assetUrl("prototype/translate.png"),
  tabCheckin: assetUrl("tab-checkin.png"),
  tabCheckinActive: assetUrl("tab-checkin-active.png"),
  tabHome: assetUrl("tab-home.png"),
  tabHomeActive: assetUrl("tab-home-active.png"),
  tabLearn: assetUrl("tab-learn.png"),
  tabLearnActive: assetUrl("tab-learn-active.png")
};

module.exports = {
  REMOTE_ASSET_BASE_URL,
  REMOTE_ASSET_VERSION,
  ASSETS,
  assetUrl
};
