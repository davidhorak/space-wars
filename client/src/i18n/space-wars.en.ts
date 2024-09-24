export default {
  title: "Space Wars",
  copyright: "Copyright © {{year}} David Horak",
  views: {
    battlefield: {
      pageTitle: "$t(title) / Battle Room",
      options: {
        title: "Options:",
        colliders: "Show Colliders",
        energy: "Show Energy",
        health: "Show Health",
        names: "Show Names",
        autoReset: "Auto Reset",
      },
      actions: {
        start: "Start",
        restart: "Restart",
        pause: "Pause",
        resume: "Resume",
        step: "Step",
      },
      scoreboard: {
        title: "Scoreboard:",
        totalRounds: "Total Rounds: {{totalRounds}}",
        overall: "Overall:",
        kills: "Kills:",
        destroyed: "Destroyed:",
      },
      log: {
        title: "Space Log:",
        damage:
          "<1>{{who}}</1> did {{damage}} damage to <2>{{whom}}</2> with {{damageType}}",
        collision: "<1>{{who}}</1> collided with <2>{{with}}</2>",
        killed: "<1>{{who}}</1> was killed by <2>{{whom}}</2>",
      },
      gameOver: {
        title: "Battle Concluded",
        winner: "Glory to: <1>{{winner}}</1>",
        score: "Score: {{score}}",
      },
    },
    pageNotFound: {
      pageTitle: "$t(title) / Page not Found",
      title: "Whoops!",
      subtitle: "This page got lost in conversation.",
      body: "Not to worry. <1>Click here</1> to go back and find pros in your area.",
    },
  },
  error: {
    title: "An error occurred",
    message: "Something went wrong. Please try again.",
    refresh: {
      button: "Refresh",
    },
  },
};
