<script>
    import Header from '../../components/Header.svelte';
    import Button from '../../components/Button.svelte';
    import { push } from 'svelte-spa-router';


</script>

<div id="pacient-hub">
    <Header return_to="/"/>
    <main id="hub-content">
        <div id="content-container">
            <div class="hub-item">
                Pending Questionarys: 0
            </div>
            <div class="hs"></div>
            <div class="hub-item">
                <Button onClick={() => push("/mri")} label="Upload MRI" isSquared={true} padding="1vh 1vw"/>
            </div>
        </div>
    </main>
</div>

<style>
    #hub-content {
        display: flex;
        height: 90vh;
        justify-content: center;
        align-items: center;
    }

    #content-container{
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        align-items: center;
        height: 20vh;
        width: 40%;
    }

    .hub-item:first-child {
        font-size: 1.7em;
    }
</style>