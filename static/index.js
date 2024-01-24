function notify(text, status) {
  UIkit.notification(
    { message: text
    , status: status
    , pos: "top-right"
    , timeout: 3000
    })
}

async function getTopics() {
  const res = await fetch("/topic")
  const topics = await res.json()
  const topicSet = new Set()
  for (const layer of topics) {
    for (const topic of layer) {
      topicSet.add(topic)
    }
  }
  return topicSet
}

async function loadTopics() {
  const topics = await getTopics()
  const select = document.getElementById("topic-select")
  for (let i = select.options.length - 1; i >= 1; i--) {
    select.remove(i)
  }
  for (const topic of topics) {
    const option = document.createElement("option")
    option.value = topic
    option.innerHTML = topic
    select.appendChild(option)
  }
  notify("Successfully loaded topics!", "success")
}

async function onSelect() {
  const select = document.getElementById("topic-select")
  const topic = document.getElementById("topic")
  topic.value = select.value
  const dataRes = await fetch("/topic/" + select.value)
  const mimetype = dataRes.headers.get("content-type")
  document.getElementById("mimetype").value = mimetype
  let data = await dataRes.text()
  if (mimetype === "application/json") {
    data = JSON.stringify(JSON.parse(data), null, "  ")
  }
  document.getElementById("data").value = data
}

async function publish() {
  let data = document.getElementById("data").value
  const mimetype = document.getElementById("mimetype").value
  if (mimetype.length == 0) {
    notify("Cannot publish with no mimetype", "danger")
    return
  }
  if (mimetype === "application/json") {
    data = JSON.stringify(JSON.parse(data), null, 0)
  }
  const topic = document.getElementById("topic").value
  console.log(topic.length)
  if (topic.length == 0) {
    notify("Cannot publish with no topic", "danger")
    return
  }
  const res = await fetch("/topic/" + topic,
    { headers:
      { "Content-Type": mimetype
      }
    , body: data
    , method: "POST"
    })
  if (res.status == 200) {
    await loadTopics()
    clearFields()
    notify("Successfully Posted Data!", "success")
  } else {
    notify("Failed to Post Data... " + await res.text(), "danger")
  }
}

function clearFields() {
  for (const id of ["topic", "mimetype", "data"]) {
    document.getElementById(id).value = ""
  }
}

function setup() {
  const yearElements = document.getElementsByTagName('Year')
  for (let ele of yearElements) {
    ele.innerHTML = new Date().getFullYear()
  }
  document.getElementById("publish").onclick = publish
  clearFields()
}

setup()
loadTopics()
