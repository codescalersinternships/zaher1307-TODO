<script lang="ts">
  import Task from "./lib/Task.svelte"
import { onMount } from "svelte";

const URL:string = "http://localhost:8010/todolist";

let list:{ id: number; todoItem: string; completed: boolean }[] = [];

let newTask:string = "";

async function fetchingData() {
  fetch(URL)
    .then(response => response.json())
    .then(data => {
      list = data;
    }).catch(error => {
    });
}

onMount(() => {
  fetchingData()
});

async function handleDelete(id: number){
  await fetch(URL + "/" + id , {
    method: "DELETE"
  });
  fetchingData();
}

async function handleSubmit(e: any){
  const doc = {
    todoItem: newTask,
    completed: false
  };

  await fetch(URL, {
    method: "POST",
    body: JSON.stringify(doc),
    headers: {'Content-Type': 'application/json'},
  })

  newTask = "";

  fetchingData();
}

</script>

<main>
  <h1>Todo App</h1>
  <div class="tasks">
  <form name='form' on:submit|preventDefault={handleSubmit}>
  <input bind:value={newTask} name='task' class="enter" type="text" placeholder="What to be done?" />
  </form>
  {#each list as t }
  <Task {handleDelete} task={t} />
  {/each}
  </div>
  </main>

<style>
  main {
    display: flex;
    align-items: center;
    flex-direction: column;
  }
h1 {
  color: #ccc;
  font-weight: 300;
  font-size: 8rem;
}
  .tasks {
    width: 30rem;
    box-shadow: -5px 5px 10px -5px rgb(23 54 71 / 50%);
  }
  .enter {
    width: 100%;
    padding: 0.5rem;
    border: none;
    font-size: 1.5rem;
    outline: none;
    border-bottom: 3px solid #ddd;
  }
  .enter::placeholder {
    color: #ccc;
    font-style: italic;
    opacity: 1;
  }
</style>
