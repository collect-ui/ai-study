Component({
  properties: {
    title: {
      type: String,
      value: ""
    },
    desc: {
      type: String,
      value: ""
    },
    icon: {
      type: String,
      value: "AI"
    },
    tag: {
      type: String,
      value: ""
    },
    actionKey: {
      type: String,
      value: ""
    },
    tone: {
      type: String,
      value: "blue"
    },
    wide: {
      type: Boolean,
      value: false
    }
  },

  methods: {
    emitAction() {
      this.triggerEvent("action", {
        key: this.properties.actionKey,
        title: this.properties.title
      });
    }
  }
});
