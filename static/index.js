const editableMimetypes = new Set(
  [ "application/json"
  , "text/plain"
  ])

const imageMimetypes = new Set(
  [ "image/jpeg"
  , "image/png"
  , "image/svg+xml"
  ])

function getBaseUrl() {
  return window.location.pathname.replace(/\/editor.*/i, "")
}

function getTopicUrl(topic) {
  return getBaseUrl() + "/topic/" + topic
}

function notify(text, status) {
  UIkit.notification(
    { message: text
    , status: status
    , pos: "top-right"
    , timeout: 3000
    })
}

async function getTopics() {
  const res = await fetch(getBaseUrl() + "/topic/")
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
  const url = getTopicUrl(select.value)
  const dataRes = await fetch(url)
  const mimetype = dataRes.headers.get("content-type")
  const mimetypeEle = document.getElementById("mimetype")
  mimetypeEle.value = mimetype
  mimetypeEle.onchange()

  if (editableMimetypes.has(mimetype)) {
    let data = await dataRes.text()
    const dataEle = document.getElementById("data")
    if (mimetype === "application/json") {
      data = JSON.stringify(JSON.parse(data), null, "  ")
    }
    dataEle.value = data
  } else if (imageMimetypes.has(mimetype)) {
    const blob = await dataRes.blob()
    console.log(blob)
    setImagePreview(blob)
  }
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
  } else {
    data = document.getElementById("file-input").files[0]
  }
  const topic = document.getElementById("topic").value
  console.log(topic.length)
  if (topic.length == 0) {
    notify("Cannot publish with no topic", "danger")
    return
  }
  const res = await fetch(getBaseUrl() + "/topic/" + topic,
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

async function deleteTopic() {
  const topic = document.getElementById("topic").value
  if (topic.length == 0) {
    notify("Cannot delete with no topic", "danger")
    return
  }
  const url = getTopicUrl(topic)
  const res = await fetch(url, { method: "DELETE" })
  if (res.status == 200) {
    await loadTopics()
    clearFields()
    notify("Successfully Deleted Topic!", "success")
  } else {
    notify("Failed to Delete Topic... " + await res.text(), "danger")
  }
}

function clearFields() {
  for (const id of ["topic", "mimetype", "data"]) {
    const ele = document.getElementById(id)
    ele.value = ""
    if (ele.onchange !== null) {
      ele.onchange()
    }
  }
}

function hideThesaurumEditControls() {
  const editControls = document.getElementsByClassName("thesaurum-edit-control")
  for (const control of editControls) {
    control.setAttribute("hidden", "hidden")
  }
}

function onMimetypeChange() {
  console.log("on mimetype change")
  hideThesaurumEditControls()
  const mimetype = document.getElementById("mimetype").value
  if (mimetype === "") return;
  if (editableMimetypes.has(mimetype)) {
    document.getElementById("data").removeAttribute("hidden")
  } else {
    document.getElementById("upload-controls").removeAttribute("hidden")
    const fileInputEle = document.getElementById("file-input")
    fileInputEle.setAttribute("accept", mimetype)
    fileInputEle.value = null
    if (imageMimetypes.has(mimetype)) {
      document.getElementById("image-preview").removeAttribute("hidden")
    }
  }
}

function onFileInputChange(e) {
  setImagePreview(e.target.files[0])
}

function setImagePreview(src) {
  const preview = document.getElementById("image-preview")
  preview.src = URL.createObjectURL(src)
  preview.onload = function() {
    URL.revokeObjectURL(preview.src)
    notify("Loaded Image!", "success")
  }
}

function setup() {
  const yearElements = document.getElementsByTagName('Year')
  for (let ele of yearElements) {
    ele.innerHTML = new Date().getFullYear()
  }
  document.getElementById("publish").onclick = publish
  document.getElementById("delete").onclick = deleteTopic
  const mimetypeEle = document.getElementById("mimetype")
  for (const type of editableMimetypes) {
    const option = document.createElement("option")
    option.value = type
    option.innerHTML = type
    mimetypeEle.appendChild(option)
  }
  for (const type of imageMimetypes) {
    const option = document.createElement("option")
    option.value = type
    option.innerHTML = type
    mimetypeEle.appendChild(option)
  }

  clearFields()
  hideThesaurumEditControls()
}

setup()
loadTopics()
