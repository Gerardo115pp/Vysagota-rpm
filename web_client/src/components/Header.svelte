<script>
    import { push } from 'svelte-spa-router';

    export let return_to = null;

    const goBack = () => {
        if (return_to) {
            push(return_to);
        }
    }
</script>

<div id="vysagota-header">
    <div id="vh-content">
        <div id="vh-label">Vysagota</div>
        <div id="vh-controls">
            <slot></slot>
            {#if return_to}
                <span on:click={goBack} id="back-btn">
                    back
                </span>
            {/if}
        </div>
    </div>
    <div class="hs" />
</div>

<style>
    #vysagota-header {
        font-family: var(--title-font);
    }

    #vh-content {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        padding: 0 2vw;
        margin: 2vh 0 ;
    }

    #vh-controls {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-end;
    }

    #vh-label {
        width: 30vw;
        font-size: 2rem;
        font-weight: normal;
    }

    #back-btn {
        cursor: pointer;
        font-size: 1.3rem;
        padding: 0 0 0 5vw;
    }

</style>
