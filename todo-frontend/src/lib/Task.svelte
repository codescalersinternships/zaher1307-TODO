<script lang="ts">
  export let task: {id:number, todoItem: string, completed: boolean};
export let handleDelete:Function;

const URL:string = "http://localhost:8010/todolist";

async function handleDone(){
  task.completed = !task.completed;


  const doc = {...task, completed : task.completed}


  await fetch(URL, {
    method: "PATCH",
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(doc)
  })
}
</script>

<div class="task">
  <button class='btn {task.completed? "btn-done":""}' on:click={handleDone} ></button>
  <div class="p--contianer">
  <p class={ task.completed? "text-done" :"" }>{task.todoItem}</p>
  </div>
  <button on:click={handleDelete(task.id)}>Delete</button>
  </div>

<style>
  .task {
    padding: 0.5rem;
    width: 100%;
    background-color: #fff;
    display: flex;
    align-items: center;
  }
  .p--contianer {
    margin-left: 1rem;
    flex: 1
  }
p {
  font-size: 1.5rem;
  width: fit-content;
  position: relative;
  cursor: default;
}

  .btn {
    width: 1rem;
    height: 1rem;
    border-radius: 50%;
    margin-left: 0.5rem;
    border: 1px solid #bbb;
    cursor: pointer;
  }
  .btn-done {
    background-color: green;
  }
  .text-done {
    color: #bbb;
  }
  .text-done::after {
    content: '';
    display: block;
    width: 100%;
    height: 50%;
    position: absolute;
    top: 0;
    left: 0;
    border-bottom: 3px solid;
  }
</style>
