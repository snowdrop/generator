package {{.PackageName}};

{{$hasDekorate := .HasDekorate}}

{{if $hasDekorate}}
import io.dekorate.kubernetes.annotation.KubernetesApplication;
import io.dekorate.openshift.annotation.OpenshiftApplication;
{{end}}
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.boot.web.support.SpringBootServletInitializer;

@SpringBootApplication
{{if $hasDekorate}}
@KubernetesApplication
@OpenshiftApplication
{{end}}
public class MvcApplication extends SpringBootServletInitializer {

	@Override
	protected SpringApplicationBuilder configure(SpringApplicationBuilder application) {
		return application.sources(MvcApplication.class);
	}

	public static void main(String[] args) throws Exception {
		SpringApplication.run(MvcApplication.class, args);
	}
}
